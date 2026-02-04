package data

import (
	"fmt"
	"path/filepath"
)

type EconomicData struct {
	GDPPerCapita          map[string]float64
	UnemploymentRate      map[string]float64
	YouthUnemploymentRate map[string]float64
	EducationExpenditure  map[string]float64
	LaborForce            map[string]float64
	InflationRate         map[string]float64
}

func LoadEconomicData(dataDir string) (*EconomicData, error) {
	e := &EconomicData{}
	var err error
	var records []StatRecord

	records, err = LoadCSV(filepath.Join(dataDir, "gdp_per_cap.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading gdp_per_cap.csv: %w", err)
	}
	e.GDPPerCapita = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "unemployment_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading unemployment_rate.csv: %w", err)
	}
	e.UnemploymentRate = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "youth_unemployment_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading youth_unemployment_rate.csv: %w", err)
	}
	e.YouthUnemploymentRate = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "education_expenditure.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading education_expenditure.csv: %w", err)
	}
	e.EducationExpenditure = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "labor_force.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading labor_force.csv: %w", err)
	}
	e.LaborForce = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "inflation_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading inflation_rate.csv: %w", err)
	}
	e.InflationRate = RecordsToMap(records)

	return e, nil
}

func (e *EconomicData) GetGDPPerCapita(slug string) float64 {
	if v, ok := e.GDPPerCapita[slug]; ok {
		return v
	}
	return 15000.0
}

func (e *EconomicData) GetUnemploymentRate(slug string) float64 {
	if v, ok := e.UnemploymentRate[slug]; ok {
		return v
	}
	return 6.0
}

func (e *EconomicData) GetYouthUnemploymentRate(slug string) float64 {
	if v, ok := e.YouthUnemploymentRate[slug]; ok {
		return v
	}
	return 15.0
}

func (e *EconomicData) GetEducationExpenditure(slug string) float64 {
	if v, ok := e.EducationExpenditure[slug]; ok {
		return v
	}
	return 4.5
}

func (e *EconomicData) GetLaborForce(slug string) float64 {
	if v, ok := e.LaborForce[slug]; ok {
		return v
	}
	return 0
}

func (e *EconomicData) GetInflationRate(slug string) float64 {
	if v, ok := e.InflationRate[slug]; ok {
		return v
	}
	return 3.0
}
