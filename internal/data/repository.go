package data

import (
	"fmt"
	"sort"
	"strings"
)

// Repository provides centralized access to all data
type Repository struct {
	Demographic *DemographicData
	Economic    *EconomicData
	Health      *HealthData
	Identity    *IdentityData
	Historical  *HistoricalData
	dataDir     string
}

// CountryStats holds all statistics for a single country
type CountryStats struct {
	Slug           string
	ISOCode        string
	Name           string
	BirthRate      float64
	DeathRate      float64
	LifeExpectancy float64
	MigrationRate  float64
	InfantMortality float64
	Population     float64
	GDPPerCapita   float64
	UnemploymentRate float64
	YouthUnemploymentRate float64
	EducationExpenditure float64
	AlcoholConsumption float64
	TobaccoUse     float64
}

// NewRepository creates a new data repository loading all data from the directory
func NewRepository(dataDir string) (*Repository, error) {
	r := &Repository{dataDir: dataDir}
	var err error

	// Load demographic data
	r.Demographic, err = LoadDemographicData(dataDir)
	if err != nil {
		return nil, fmt.Errorf("loading demographic data: %w", err)
	}

	// Load economic data
	r.Economic, err = LoadEconomicData(dataDir)
	if err != nil {
		return nil, fmt.Errorf("loading economic data: %w", err)
	}

	// Load health data
	r.Health, err = LoadHealthData(dataDir)
	if err != nil {
		return nil, fmt.Errorf("loading health data: %w", err)
	}

	// Load identity data
	r.Identity, err = LoadIdentityData(dataDir)
	if err != nil {
		return nil, fmt.Errorf("loading identity data: %w", err)
	}

	// Load historical data
	r.Historical, err = LoadHistoricalData(dataDir)
	if err != nil {
		return nil, fmt.Errorf("loading historical data: %w", err)
	}

	return r, nil
}

// GetCountryStats returns all statistics for a country by slug
func (r *Repository) GetCountryStats(slug string) *CountryStats {
	isoCode := r.Identity.GetISOCodeFromSlug(slug)

	return &CountryStats{
		Slug:                  slug,
		ISOCode:               isoCode,
		Name:                  r.Identity.GetCountryName(isoCode),
		BirthRate:             r.Demographic.GetBirthRate(slug),
		DeathRate:             r.Demographic.GetDeathRate(slug),
		LifeExpectancy:        r.Demographic.GetLifeExpectancy(slug),
		MigrationRate:         r.Demographic.GetMigrationRate(slug),
		InfantMortality:       r.Demographic.GetInfantMortality(slug),
		Population:            r.Demographic.GetPopulation(slug),
		GDPPerCapita:          r.Economic.GetGDPPerCapita(slug),
		UnemploymentRate:      r.Economic.GetUnemploymentRate(slug),
		YouthUnemploymentRate: r.Economic.GetYouthUnemploymentRate(slug),
		EducationExpenditure:  r.Economic.GetEducationExpenditure(slug),
		AlcoholConsumption:    r.Health.GetAlcoholConsumption(slug),
		TobaccoUse:            r.Health.GetTobaccoUse(slug),
	}
}

// GetAvailableCountrySlugs returns all country slugs that have demographic data
func (r *Repository) GetAvailableCountrySlugs() []string {
	slugs := make([]string, 0, len(r.Demographic.BirthRates))
	for slug := range r.Demographic.BirthRates {
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)
	return slugs
}

// GetCountriesWithNames returns countries that have both demographic and name data
func (r *Repository) GetCountriesWithNames() []string {
	var countries []string

	for slug := range r.Demographic.BirthRates {
		isoCode := r.Identity.GetISOCodeFromSlug(slug)
		if isoCode != "" {
			forenames := r.Identity.GetForenames(isoCode)
			surnames := r.Identity.GetSurnames(isoCode)
			if len(forenames) > 0 && len(surnames) > 0 {
				countries = append(countries, slug)
			}
		}
	}

	sort.Strings(countries)
	return countries
}

// ValidateCountry checks if a country slug is valid and has sufficient data
func (r *Repository) ValidateCountry(slug string) error {
	slug = strings.ToLower(slug)

	// Check demographic data
	if _, ok := r.Demographic.BirthRates[slug]; !ok {
		return fmt.Errorf("country '%s' not found in demographic data", slug)
	}

	// Check for name data
	isoCode := r.Identity.GetISOCodeFromSlug(slug)
	if isoCode == "" {
		return fmt.Errorf("country '%s' has no ISO code mapping", slug)
	}

	forenames := r.Identity.GetForenames(isoCode)
	if len(forenames) == 0 {
		return fmt.Errorf("country '%s' has no forename data", slug)
	}

	surnames := r.Identity.GetSurnames(isoCode)
	if len(surnames) == 0 {
		return fmt.Errorf("country '%s' has no surname data", slug)
	}

	return nil
}

// GetISOCodeForSlug returns the ISO code for a country slug
func (r *Repository) GetISOCodeForSlug(slug string) string {
	return r.Identity.GetISOCodeFromSlug(slug)
}

// GetForenamesByGender returns forenames for a country and gender
func (r *Repository) GetForenamesByGender(slug, gender string) []NameRecord {
	isoCode := r.Identity.GetISOCodeFromSlug(slug)
	if isoCode == "" {
		return nil
	}
	return r.Identity.GetForenamesByGender(isoCode, gender)
}

// GetSurnames returns surnames for a country
func (r *Repository) GetSurnames(slug string) []SurnameRecord {
	isoCode := r.Identity.GetISOCodeFromSlug(slug)
	if isoCode == "" {
		return nil
	}
	return r.Identity.GetSurnames(isoCode)
}

func (r *Repository) GetGDPPerCapita(slug string) float64 {
	return r.Economic.GetGDPPerCapita(slug)
}

func (r *Repository) GetUnderweightU5(slug string) float64 {
	return r.Health.GetUnderweightU5(slug)
}

// GetFertilityRate returns the total fertility rate for a country and year
func (r *Repository) GetFertilityRate(slug string, year int) float64 {
	iso3 := GetISO3FromSlug(slug)
	if iso3 == "" {
		return 2.1 // Replacement rate default
	}
	return r.Historical.FertilityRate.GetValueOrDefault(iso3, year, 2.1)
}

// GetMarriageAgeWomen returns the mean marriage age for women
func (r *Repository) GetMarriageAgeWomen(slug string, year int) float64 {
	iso3 := GetISO3FromSlug(slug)
	if iso3 == "" {
		return 25.0 // Default
	}
	return r.Historical.MarriageAgeWomen.GetValueOrDefault(iso3, year, 25.0)
}

// GetDivorceRate returns divorces per 1000 people
func (r *Repository) GetDivorceRate(slug string, year int) float64 {
	iso3 := GetISO3FromSlug(slug)
	if iso3 == "" {
		return 2.0 // Default
	}
	return r.Historical.DivorceRate.GetValueOrDefault(iso3, year, 2.0)
}

// GetYouthMortality returns under-15 mortality rate
func (r *Repository) GetYouthMortality(slug string, year int) float64 {
	iso3 := GetISO3FromSlug(slug)
	if iso3 == "" {
		return 5.0 // Default
	}
	return r.Historical.YouthMortality.GetValueOrDefault(iso3, year, 5.0)
}

// GetBirthsOutsideMarriage returns share of births outside marriage (%)
func (r *Repository) GetBirthsOutsideMarriage(slug string, year int) float64 {
	iso3 := GetISO3FromSlug(slug)
	if iso3 == "" {
		return 20.0 // Default
	}
	return r.Historical.BirthsOutsideMarriage.GetValueOrDefault(iso3, year, 20.0)
}

func (r *Repository) GetUrbanShare(slug string, year int) float64 {
	iso3 := GetISO3FromSlug(slug)
	if iso3 == "" || r.Historical.UrbanPopulationShare == nil {
		return 0.5
	}
	value, ok := r.Historical.UrbanPopulationShare.GetValue(iso3, year)
	if !ok {
		return 0.5
	}
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

// GetMarriageRate returns marriages per 1000 people
func (r *Repository) GetMarriageRate(slug string, year int) float64 {
	iso3 := GetISO3FromSlug(slug)
	if iso3 == "" {
		return 5.0 // Default
	}
	return r.Historical.MarriageRate.GetValueOrDefault(iso3, year, 5.0)
}

// GetSingleParentShare returns share of single-parent households (%)
func (r *Repository) GetSingleParentShare(slug string, year int) float64 {
	iso3 := GetISO3FromSlug(slug)
	if iso3 == "" {
		return 10.0 // Default
	}
	return r.Historical.SingleParentShare.GetValueOrDefault(iso3, year, 10.0)
}
