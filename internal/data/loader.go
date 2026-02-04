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

type StatRecord struct {
	Name     string
	Slug     string
	Value    float64
	RawValue string
	Year     int
	Ranking  int
	Region   string
}

func LoadCSV(filepath string) ([]StatRecord, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filepath, err)
	}

	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV %s: %w", filepath, err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file %s has no data rows", filepath)
	}

	var result []StatRecord
	for i, row := range records[1:] {
		if len(row) < 6 {
			continue
		}

		record := StatRecord{
			Name:     strings.Trim(row[0], "\""),
			Slug:     strings.Trim(row[1], "\""),
			RawValue: strings.Trim(row[2], "\""),
			Region:   strings.Trim(row[len(row)-1], "\""),
		}

		record.Value, _ = ParseValue(record.RawValue)

		yearIdx := 3
		if len(row) > 6 {
			yearIdx = len(row) - 3
		}
		if yearIdx < len(row) {
			yearStr := strings.Trim(row[yearIdx], "\"")
			record.Year, _ = strconv.Atoi(yearStr)
		}

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

func ParseValue(raw string) (float64, error) {

	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "$")
	s = strings.TrimSuffix(s, "%")

	s = strings.ReplaceAll(s, ",", "")

	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = "-" + s[1:len(s)-1]
	}

	if s == "" || s == "N/A" || s == "-" {
		return 0, fmt.Errorf("empty or invalid value")
	}

	return strconv.ParseFloat(s, 64)
}

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

func RecordsToMap(records []StatRecord) map[string]float64 {
	result := make(map[string]float64)
	for _, r := range records {
		result[r.Slug] = r.Value
	}
	return result
}

func RecordsToFullMap(records []StatRecord) map[string]StatRecord {
	result := make(map[string]StatRecord)
	for _, r := range records {
		result[r.Slug] = r
	}
	return result
}
