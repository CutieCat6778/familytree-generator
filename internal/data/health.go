package data

import (
	"fmt"
	"path/filepath"
)

type HealthData struct {
	AlcoholConsumption map[string]float64
	TobaccoUse         map[string]float64
	UnderweightU5      map[string]float64
}

func LoadHealthData(dataDir string) (*HealthData, error) {
	h := &HealthData{}
	var err error
	var records []StatRecord

	records, err = LoadCSV(filepath.Join(dataDir, "alcohol.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading alcohol.csv: %w", err)
	}
	h.AlcoholConsumption = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "tobacco_use.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading tobacco_use.csv: %w", err)
	}
	h.TobaccoUse = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "underweight_u5.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading underweight_u5.csv: %w", err)
	}
	h.UnderweightU5 = RecordsToMap(records)

	return h, nil
}

func (h *HealthData) GetAlcoholConsumption(slug string) float64 {
	if v, ok := h.AlcoholConsumption[slug]; ok {
		return v
	}
	return 6.0
}

func (h *HealthData) GetTobaccoUse(slug string) float64 {
	if v, ok := h.TobaccoUse[slug]; ok {
		return v
	}
	return 20.0
}

func (h *HealthData) GetUnderweightU5(slug string) float64 {
	if v, ok := h.UnderweightU5[slug]; ok {
		return v
	}
	return 15.0
}
