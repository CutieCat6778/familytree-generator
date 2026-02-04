package data

import (
	"fmt"
	"path/filepath"
)

// EconomicData holds all economic statistics
type EconomicData struct {
	GDPPerCapita          map[string]float64 // slug -> USD
	UnemploymentRate      map[string]float64 // slug -> percentage
	YouthUnemploymentRate map[string]float64 // slug -> percentage
	EducationExpenditure  map[string]float64 // slug -> % of GDP
	LaborForce            map[string]float64 // slug -> count
	InflationRate         map[string]float64 // slug -> percentage
}

// LoadEconomicData loads all economic CSV files from the data directory
func LoadEconomicData(dataDir string) (*EconomicData, error) {
	e := &EconomicData{}
	var err error
	var records []StatRecord

	// GDP per capita
	records, err = LoadCSV(filepath.Join(dataDir, "gdp_per_cap.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading gdp_per_cap.csv: %w", err)
	}
	e.GDPPerCapita = RecordsToMap(records)

	// Unemployment rate
	records, err = LoadCSV(filepath.Join(dataDir, "unemployment_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading unemployment_rate.csv: %w", err)
	}
	e.UnemploymentRate = RecordsToMap(records)

	// Youth unemployment rate
	records, err = LoadCSV(filepath.Join(dataDir, "youth_unemployment_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading youth_unemployment_rate.csv: %w", err)
	}
	e.YouthUnemploymentRate = RecordsToMap(records)

	// Education expenditure
	records, err = LoadCSV(filepath.Join(dataDir, "education_expenditure.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading education_expenditure.csv: %w", err)
	}
	e.EducationExpenditure = RecordsToMap(records)

	// Labor force
	records, err = LoadCSV(filepath.Join(dataDir, "labor_force.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading labor_force.csv: %w", err)
	}
	e.LaborForce = RecordsToMap(records)

	// Inflation rate
	records, err = LoadCSV(filepath.Join(dataDir, "inflation_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading inflation_rate.csv: %w", err)
	}
	e.InflationRate = RecordsToMap(records)

	return e, nil
}

// GetGDPPerCapita returns the GDP per capita for a country, or a default value
func (e *EconomicData) GetGDPPerCapita(slug string) float64 {
	if v, ok := e.GDPPerCapita[slug]; ok {
		return v
	}
	return 15000.0 // World average approximation
}

// GetUnemploymentRate returns the unemployment rate for a country, or a default value
func (e *EconomicData) GetUnemploymentRate(slug string) float64 {
	if v, ok := e.UnemploymentRate[slug]; ok {
		return v
	}
	return 6.0 // World average approximation
}

// GetYouthUnemploymentRate returns the youth unemployment rate for a country, or a default value
func (e *EconomicData) GetYouthUnemploymentRate(slug string) float64 {
	if v, ok := e.YouthUnemploymentRate[slug]; ok {
		return v
	}
	return 15.0 // World average approximation
}

// GetEducationExpenditure returns the education expenditure for a country, or a default value
func (e *EconomicData) GetEducationExpenditure(slug string) float64 {
	if v, ok := e.EducationExpenditure[slug]; ok {
		return v
	}
	return 4.5 // World average approximation
}

// GetLaborForce returns the labor force for a country, or 0
func (e *EconomicData) GetLaborForce(slug string) float64 {
	if v, ok := e.LaborForce[slug]; ok {
		return v
	}
	return 0
}

// GetInflationRate returns the inflation rate for a country, or a default value
func (e *EconomicData) GetInflationRate(slug string) float64 {
	if v, ok := e.InflationRate[slug]; ok {
		return v
	}
	return 3.0 // Target rate approximation
}
