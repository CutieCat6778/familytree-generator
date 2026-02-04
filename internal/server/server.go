package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/generator"
	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/internal/output"
)

type Server struct {
	repo        *data.Repository
	addr        string
	webDir      string
	rateLimiter *RateLimiter
}

func NewServer(repo *data.Repository, addr string, webDir string) *Server {
	return &Server{
		repo:        repo,
		addr:        addr,
		webDir:      webDir,
		rateLimiter: NewRateLimiter(10, time.Minute),
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/generate", s.corsMiddleware(s.handleGenerate))
	mux.HandleFunc("/api/countries", s.corsMiddleware(s.handleCountries))
	mux.HandleFunc("/api/country/", s.corsMiddleware(s.handleCountryStats))
	mux.HandleFunc("/api/health", s.corsMiddleware(s.handleHealth))

	if s.webDir != "" {
		fs := http.FileServer(http.Dir(s.webDir))
		mux.Handle("/", fs)
	}

	log.Printf("Server starting on %s", s.addr)
	log.Printf("API endpoints:")
	log.Printf("  GET  /api/health - Health check")
	log.Printf("  GET  /api/countries - List available countries")
	log.Printf("  GET  /api/country/{slug} - Get country statistics")
	log.Printf("  POST /api/generate - Generate a family tree")

	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSuffix(r.Header.Get("Origin"), "/")
		if origin != "" && !isAllowedOrigin(origin) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if reqHeaders := r.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
			w.Header().Add("Vary", "Access-Control-Request-Headers")
		} else {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func isAllowedOrigin(origin string) bool {
	if strings.HasPrefix(origin, "https://familytree.thinis.de") {
		return true
	}
	if strings.HasPrefix(origin, "http://localhost") {
		return true
	}
	if strings.HasPrefix(origin, "http://127.0.0.1") {
		return true
	}
	return false
}

type GenerateRequest struct {
	Country            string `json:"country"`
	Generations        int    `json:"generations"`
	Seed               int64  `json:"seed"`
	StartYear          int    `json:"start_year"`
	Gender             string `json:"gender"`
	IncludeExtended    bool   `json:"include_extended"`
	LifeExpectancyMode string `json:"life_expectancy_mode"`
}

type GenerateResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message,omitempty"`
	Tree    *output.VisualizationData `json:"tree,omitempty"`
	Stats   *TreeStats                `json:"stats,omitempty"`
}

type TreeStats struct {
	TotalPersons   int     `json:"total_persons"`
	TotalFamilies  int     `json:"total_families"`
	LivingPersons  int     `json:"living_persons"`
	AverageAge     float64 `json:"average_age"`
	OldestPerson   int     `json:"oldest_person_age"`
	TotalChildren  int     `json:"total_children"`
	DivorceCount   int     `json:"divorce_count"`
	GenerationTime string  `json:"generation_time"`
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.rateLimiter != nil && !s.rateLimiter.Allow(clientKey(r)) {
		s.jsonError(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Country == "" {
		req.Country = "germany"
	}
	if req.Generations < 1 {
		req.Generations = 3
	}
	if req.Generations > 10 {
		req.Generations = 10
	}
	if req.Seed == 0 {
		req.Seed = time.Now().UnixNano()
	}
	if req.StartYear == 0 {
		req.StartYear = 1970
	}
	if req.LifeExpectancyMode == "" {
		req.LifeExpectancyMode = string(generator.LifeExpectancyTotal)
	}

	if err := s.repo.ValidateCountry(req.Country); err != nil {
		s.jsonError(w, "Invalid country: "+err.Error(), http.StatusBadRequest)
		return
	}

	var gender model.Gender
	switch strings.ToUpper(req.Gender) {
	case "M", "MALE":
		gender = model.Male
	case "F", "FEMALE":
		gender = model.Female
	}

	config := generator.Config{
		Country:            req.Country,
		Generations:        req.Generations,
		Seed:               req.Seed,
		StartYear:          req.StartYear,
		RootGender:         gender,
		IncludeExtended:    req.IncludeExtended,
		LifeExpectancyMode: generator.ParseLifeExpectancyMode(req.LifeExpectancyMode),
	}

	startTime := time.Now()
	engine := generator.NewEngine(config, s.repo)
	tree, err := engine.Generate()
	if err != nil {
		s.jsonError(w, "Generation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	generationTime := time.Since(startTime)

	vizData := output.TreeToVisualizationData(tree)

	divorceCount := 0
	for _, family := range tree.GetAllFamilies() {
		if family.IsDivorced() {
			divorceCount++
		}
	}

	stats := &TreeStats{
		TotalPersons:   vizData.Stats.TotalPersons,
		TotalFamilies:  vizData.Stats.TotalFamilies,
		LivingPersons:  vizData.Stats.LivingPersons,
		AverageAge:     vizData.Stats.AverageAge,
		OldestPerson:   vizData.Stats.OldestPerson,
		TotalChildren:  vizData.Stats.TotalChildren,
		DivorceCount:   divorceCount,
		GenerationTime: generationTime.String(),
	}

	response := GenerateResponse{
		Success: true,
		Tree:    vizData,
		Stats:   stats,
	}

	s.jsonResponse(w, response)
}

type CountryInfo struct {
	Slug           string  `json:"slug"`
	Name           string  `json:"name"`
	ISOCode        string  `json:"iso_code"`
	HasNameData    bool    `json:"has_name_data"`
	Population     float64 `json:"population,omitempty"`
	LifeExpectancy float64 `json:"life_expectancy,omitempty"`
}

func (s *Server) handleCountries(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		s.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	countries := s.repo.GetCountriesWithNames()
	result := make([]CountryInfo, 0, len(countries))

	for _, slug := range countries {
		stats := s.repo.GetCountryStats(slug)
		info := CountryInfo{
			Slug:           slug,
			Name:           stats.Name,
			ISOCode:        stats.ISOCode,
			HasNameData:    true,
			Population:     stats.Population,
			LifeExpectancy: stats.LifeExpectancy,
		}
		result = append(result, info)
	}

	s.jsonResponse(w, map[string]interface{}{
		"countries": result,
		"count":     len(result),
	})
}

func (s *Server) handleCountryStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		s.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/country/")
	if slug == "" {
		s.jsonError(w, "Country slug required", http.StatusBadRequest)
		return
	}

	if err := s.repo.ValidateCountry(slug); err != nil {
		s.jsonError(w, "Country not found: "+err.Error(), http.StatusNotFound)
		return
	}

	stats := s.repo.GetCountryStats(slug)

	currentYear := time.Now().Year()
	historicalStats := map[string]interface{}{
		"fertility_rate":          s.repo.GetFertilityRate(slug, currentYear),
		"marriage_age_women":      s.repo.GetMarriageAgeWomen(slug, currentYear),
		"divorce_rate":            s.repo.GetDivorceRate(slug, currentYear),
		"youth_mortality":         s.repo.GetYouthMortality(slug, currentYear),
		"births_outside_marriage": s.repo.GetBirthsOutsideMarriage(slug, currentYear),
		"marriage_rate":           s.repo.GetMarriageRate(slug, currentYear),
		"single_parent_share":     s.repo.GetSingleParentShare(slug, currentYear),
	}

	s.jsonResponse(w, map[string]interface{}{
		"slug":       stats.Slug,
		"name":       stats.Name,
		"iso_code":   stats.ISOCode,
		"current":    stats,
		"historical": historicalStats,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.jsonResponse(w, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (s *Server) jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func ParsePort(port string) (int, error) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("invalid port: %s", port)
	}
	if p < 1 || p > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535")
	}
	return p, nil
}
