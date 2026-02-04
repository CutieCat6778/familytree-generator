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

type NameRecord struct {
	Country       string
	CountryGroup  int
	Region        string
	Year          int
	Index         int
	Gender        string
	LocalizedName string
	RomanizedName string
}

type SurnameRecord struct {
	Country       string
	Rank          int
	LocalizedName string
	RomanizedName string
	Percentage    float64
	Count         int
}

type IdentityData struct {
	Forenames    map[string][]NameRecord
	Surnames     map[string][]SurnameRecord
	CountryCodes map[string]string

	NameToCode map[string]string

	SlugToCode map[string]string
}

func LoadIdentityData(dataDir string) (*IdentityData, error) {
	id := &IdentityData{
		Forenames:    make(map[string][]NameRecord),
		Surnames:     make(map[string][]SurnameRecord),
		CountryCodes: make(map[string]string),
		NameToCode:   make(map[string]string),
		SlugToCode:   make(map[string]string),
	}

	if err := LoadJSON(filepath.Join(dataDir, "countries-code.json"), &id.CountryCodes); err != nil {
		return nil, fmt.Errorf("loading countries-code.json: %w", err)
	}

	for code, name := range id.CountryCodes {
		nameLower := strings.ToLower(name)
		id.NameToCode[nameLower] = code

		slug := toSlug(name)
		id.SlugToCode[slug] = code
	}

	if err := id.loadForenames(filepath.Join(dataDir, "forenames.csv")); err != nil {
		return nil, fmt.Errorf("loading forenames.csv: %w", err)
	}

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

	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

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

	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, row := range records[1:] {
		if len(row) < 6 {
			continue
		}

		rank, _ := strconv.Atoi(row[1])

		percentage, err := strconv.ParseFloat(strings.TrimSpace(row[7]), 64)
		if err != nil {
			percentage = 0.0
		}

		count, err := strconv.Atoi(strings.TrimSpace(row[6]))
		if err != nil {
			count = 0
		}

		record := SurnameRecord{
			Country:       strings.TrimSpace(row[0]),
			Rank:          rank,
			LocalizedName: strings.TrimSpace(row[4]),
			RomanizedName: strings.TrimSpace(row[5]),
			Percentage:    percentage,
			Count:         count,
		}

		if record.RomanizedName == "" {
			record.RomanizedName = record.LocalizedName
		}

		id.Surnames[record.Country] = append(id.Surnames[record.Country], record)
	}

	return nil
}

func (id *IdentityData) GetForenames(isoCode string) []NameRecord {
	return id.Forenames[isoCode]
}

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

func (id *IdentityData) GetSurnames(isoCode string) []SurnameRecord {
	return id.Surnames[isoCode]
}

func (id *IdentityData) GetCountryName(isoCode string) string {
	return id.CountryCodes[isoCode]
}

func (id *IdentityData) GetISOCode(name string) string {
	return id.NameToCode[strings.ToLower(name)]
}

func (id *IdentityData) GetISOCodeFromSlug(slug string) string {
	return id.SlugToCode[slug]
}

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
