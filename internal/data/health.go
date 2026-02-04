package data

import (
	"fmt"
	"path/filepath"
)

// HealthData holds all health-related statistics
type HealthData struct {
	AlcoholConsumption map[string]float64 // slug -> liters per capita
	TobaccoUse         map[string]float64 // slug -> percentage of population
	UnderweightU5      map[string]float64 // slug -> percentage of children under 5
}

// LoadHealthData loads all health CSV files from the data directory
func LoadHealthData(dataDir string) (*HealthData, error) {
	h := &HealthData{}
	var err error
	var records []StatRecord

	// Alcohol consumption
	records, err = LoadCSV(filepath.Join(dataDir, "alcohol.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading alcohol.csv: %w", err)
	}
	h.AlcoholConsumption = RecordsToMap(records)

	// Tobacco use
	records, err = LoadCSV(filepath.Join(dataDir, "tobacco_use.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading tobacco_use.csv: %w", err)
	}
	h.TobaccoUse = RecordsToMap(records)

	// Underweight children under 5
	records, err = LoadCSV(filepath.Join(dataDir, "underweight_u5.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading underweight_u5.csv: %w", err)
	}
	h.UnderweightU5 = RecordsToMap(records)

	return h, nil
}

// GetAlcoholConsumption returns the alcohol consumption for a country, or a default value
func (h *HealthData) GetAlcoholConsumption(slug string) float64 {
	if v, ok := h.AlcoholConsumption[slug]; ok {
		return v
	}
	return 6.0 // World average approximation
}

// GetTobaccoUse returns the tobacco use percentage for a country, or a default value
func (h *HealthData) GetTobaccoUse(slug string) float64 {
	if v, ok := h.TobaccoUse[slug]; ok {
		return v
	}
	return 20.0 // World average approximation
}

// GetUnderweightU5 returns the underweight children percentage for a country, or a default value
func (h *HealthData) GetUnderweightU5(slug string) float64 {
	if v, ok := h.UnderweightU5[slug]; ok {
		return v
	}
	return 15.0 // World average approximation
}
