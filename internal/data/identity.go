package data

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// NameRecord represents a forename from the forenames.csv
type NameRecord struct {
	Country       string // ISO 2-letter code
	CountryGroup  int
	Region        string
	Year          int
	Index         int    // Popularity rank
	Gender        string // "M" or "F"
	LocalizedName string
	RomanizedName string
}

// SurnameRecord represents a surname from the surnames.csv
type SurnameRecord struct {
	Country       string // ISO 2-letter code
	Rank          int
	LocalizedName string
	RomanizedName string
}

// IdentityData holds all name-related data
type IdentityData struct {
	Forenames    map[string][]NameRecord    // ISO code -> names list
	Surnames     map[string][]SurnameRecord // ISO code -> surnames list
	CountryCodes map[string]string          // ISO code -> country name
	// Reverse lookup: country name (lowercase) -> ISO code
	NameToCode map[string]string
	// Slug to ISO code mapping
	SlugToCode map[string]string
}

// LoadIdentityData loads all identity-related data files
func LoadIdentityData(dataDir string) (*IdentityData, error) {
	id := &IdentityData{
		Forenames:    make(map[string][]NameRecord),
		Surnames:     make(map[string][]SurnameRecord),
		CountryCodes: make(map[string]string),
		NameToCode:   make(map[string]string),
		SlugToCode:   make(map[string]string),
	}

	// Load country codes first
	if err := LoadJSON(filepath.Join(dataDir, "countries-code.json"), &id.CountryCodes); err != nil {
		return nil, fmt.Errorf("loading countries-code.json: %w", err)
	}

	// Build reverse lookups
	for code, name := range id.CountryCodes {
		nameLower := strings.ToLower(name)
		id.NameToCode[nameLower] = code

		// Create slug from name
		slug := toSlug(name)
		id.SlugToCode[slug] = code
	}

	// Load forenames
	if err := id.loadForenames(filepath.Join(dataDir, "forenames.csv")); err != nil {
		return nil, fmt.Errorf("loading forenames.csv: %w", err)
	}

	// Load surnames
	if err := id.loadSurnames(filepath.Join(dataDir, "surnames.csv")); err != nil {
		return nil, fmt.Errorf("loading surnames.csv: %w", err)
	}

	return id, nil
}

func (id *IdentityData) loadForenames(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Remove BOM
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header: Country,Country Group,Region,Population,Note,Year,Romanization,Index,Name Group,Gender,Localized Name,Romanized Name
	for _, row := range records[1:] {
		if len(row) < 12 {
			continue
		}

		countryGroup, _ := strconv.Atoi(row[1])
		year, _ := strconv.Atoi(row[5])
		index, _ := strconv.Atoi(row[7])

		record := NameRecord{
			Country:       strings.TrimSpace(row[0]),
			CountryGroup:  countryGroup,
			Region:        strings.TrimSpace(row[2]),
			Year:          year,
			Index:         index,
			Gender:        strings.TrimSpace(row[9]),
			LocalizedName: strings.TrimSpace(row[10]),
			RomanizedName: strings.TrimSpace(row[11]),
		}

		// Use romanized name if available, otherwise localized
		if record.RomanizedName == "" {
			record.RomanizedName = record.LocalizedName
		}

		id.Forenames[record.Country] = append(id.Forenames[record.Country], record)
	}

	return nil
}

func (id *IdentityData) loadSurnames(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Remove BOM
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header: Country,Rank,Index,Name Group,Localized Name,Romanized Name,Count,Percent
	for _, row := range records[1:] {
		if len(row) < 6 {
			continue
		}

		rank, _ := strconv.Atoi(row[1])

		record := SurnameRecord{
			Country:       strings.TrimSpace(row[0]),
			Rank:          rank,
			LocalizedName: strings.TrimSpace(row[4]),
			RomanizedName: strings.TrimSpace(row[5]),
		}

		// Use romanized name if available, otherwise localized
		if record.RomanizedName == "" {
			record.RomanizedName = record.LocalizedName
		}

		id.Surnames[record.Country] = append(id.Surnames[record.Country], record)
	}

	return nil
}

// GetForenames returns forenames for a country by ISO code
func (id *IdentityData) GetForenames(isoCode string) []NameRecord {
	return id.Forenames[isoCode]
}

// GetForenamesByGender returns forenames filtered by gender
func (id *IdentityData) GetForenamesByGender(isoCode, gender string) []NameRecord {
	all := id.Forenames[isoCode]
	var filtered []NameRecord
	for _, n := range all {
		if n.Gender == gender {
			filtered = append(filtered, n)
		}
	}
	return filtered
}

// GetSurnames returns surnames for a country by ISO code
func (id *IdentityData) GetSurnames(isoCode string) []SurnameRecord {
	return id.Surnames[isoCode]
}

// GetCountryName returns the country name for an ISO code
func (id *IdentityData) GetCountryName(isoCode string) string {
	return id.CountryCodes[isoCode]
}

// GetISOCode returns the ISO code for a country name
func (id *IdentityData) GetISOCode(name string) string {
	return id.NameToCode[strings.ToLower(name)]
}

// GetISOCodeFromSlug returns the ISO code for a slug
func (id *IdentityData) GetISOCodeFromSlug(slug string) string {
	return id.SlugToCode[slug]
}

// GetAvailableCountries returns all countries with name data
func (id *IdentityData) GetAvailableCountries() []string {
	countries := make([]string, 0)
	seen := make(map[string]bool)

	for code := range id.Forenames {
		if name := id.CountryCodes[code]; name != "" && !seen[code] {
			seen[code] = true
			countries = append(countries, name)
		}
	}

	return countries
}

// toSlug converts a country name to a URL-friendly slug
func toSlug(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	return s
}
