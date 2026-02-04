package data

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// StatRecord represents a single row from demographic CSV files
type StatRecord struct {
	Name     string  // Country name
	Slug     string  // URL-friendly name (e.g., "united-states")
	Value    float64 // The statistic value
	RawValue string  // Original string value
	Year     int     // Year of information
	Ranking  int     // Global ranking
	Region   string  // Geographic region
}

// LoadCSV loads a demographic CSV file and returns records
func LoadCSV(filepath string) ([]StatRecord, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filepath, err)
	}

	// Remove BOM if present
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // Allow variable field counts (header may have comma in column name)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV %s: %w", filepath, err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file %s has no data rows", filepath)
	}

	// Skip header row
	var result []StatRecord
	for i, row := range records[1:] {
		if len(row) < 6 {
			continue // Skip malformed rows
		}

		record := StatRecord{
			Name:     strings.Trim(row[0], "\""),
			Slug:     strings.Trim(row[1], "\""),
			RawValue: strings.Trim(row[2], "\""),
			Region:   strings.Trim(row[len(row)-1], "\""), // Last column is always region
		}

		// Parse value
		record.Value, _ = ParseValue(record.RawValue)

		// Parse year - it's at index 3 for 6-field rows
		yearIdx := 3
		if len(row) > 6 {
			yearIdx = len(row) - 3 // Adjust for rows with more fields
		}
		if yearIdx < len(row) {
			yearStr := strings.Trim(row[yearIdx], "\"")
			record.Year, _ = strconv.Atoi(yearStr)
		}

		// Parse ranking - it's at index 4 for 6-field rows
		rankIdx := 4
		if len(row) > 6 {
			rankIdx = len(row) - 2
		}
		if rankIdx < len(row) {
			rankStr := strings.Trim(row[rankIdx], "\"")
			record.Ranking, _ = strconv.Atoi(rankStr)
		}
		if record.Ranking == 0 {
			record.Ranking = i + 1
		}

		result = append(result, record)
	}

	return result, nil
}

// ParseValue parses a string value that may contain:
// - Numbers: "46.6"
// - Percentages: "16.4%"
// - Currency: "$270,100"
// - Large numbers: "1,416,043,270"
func ParseValue(raw string) (float64, error) {
	// Remove common prefixes/suffixes
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "$")
	s = strings.TrimSuffix(s, "%")

	// Remove commas
	s = strings.ReplaceAll(s, ",", "")

	// Handle parentheses for negative numbers
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = "-" + s[1:len(s)-1]
	}

	if s == "" || s == "N/A" || s == "-" {
		return 0, fmt.Errorf("empty or invalid value")
	}

	return strconv.ParseFloat(s, 64)
}

// LoadJSON loads a JSON file into the target interface
func LoadJSON(filepath string, target interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", filepath, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("parsing JSON %s: %w", filepath, err)
	}

	return nil
}

// RecordsToMap converts a slice of StatRecords to a map keyed by slug
func RecordsToMap(records []StatRecord) map[string]float64 {
	result := make(map[string]float64)
	for _, r := range records {
		result[r.Slug] = r.Value
	}
	return result
}

// RecordsToFullMap converts records to a map with full record info
func RecordsToFullMap(records []StatRecord) map[string]StatRecord {
	result := make(map[string]StatRecord)
	for _, r := range records {
		result[r.Slug] = r
	}
	return result
}
