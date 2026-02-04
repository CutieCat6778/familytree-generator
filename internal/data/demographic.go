package data

import (
	"fmt"
	"path/filepath"
)

type DemographicData struct {
	BirthRates      map[string]float64
	DeathRates      map[string]float64
	LifeExpectancy  map[string]float64
	MigrationRates  map[string]float64
	InfantMortality map[string]float64
	Population      map[string]float64
}

func LoadDemographicData(dataDir string) (*DemographicData, error) {
	d := &DemographicData{}
	var err error

	records, err := LoadCSV(filepath.Join(dataDir, "birth_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading birth_rate.csv: %w", err)
	}
	d.BirthRates = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "death_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading death_rate.csv: %w", err)
	}
	d.DeathRates = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "life_exp_at_birth.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading life_exp_at_birth.csv: %w", err)
	}
	d.LifeExpectancy = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "migration_rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading migration_rate.csv: %w", err)
	}
	d.MigrationRates = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "imr.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading imr.csv: %w", err)
	}
	d.InfantMortality = RecordsToMap(records)

	records, err = LoadCSV(filepath.Join(dataDir, "population.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading population.csv: %w", err)
	}
	d.Population = RecordsToMap(records)

	return d, nil
}

func (d *DemographicData) GetBirthRate(slug string) float64 {
	if v, ok := d.BirthRates[slug]; ok {
		return v
	}
	return 12.0
}

func (d *DemographicData) GetDeathRate(slug string) float64 {
	if v, ok := d.DeathRates[slug]; ok {
		return v
	}
	return 8.0
}

func (d *DemographicData) GetLifeExpectancy(slug string) float64 {
	if v, ok := d.LifeExpectancy[slug]; ok {
		return v
	}
	return 72.0
}

func (d *DemographicData) GetMigrationRate(slug string) float64 {
	if v, ok := d.MigrationRates[slug]; ok {
		return v
	}
	return 0.0
}

func (d *DemographicData) GetInfantMortality(slug string) float64 {
	if v, ok := d.InfantMortality[slug]; ok {
		return v
	}
	return 30.0
}

func (d *DemographicData) GetPopulation(slug string) float64 {
	if v, ok := d.Population[slug]; ok {
		return v
	}
	return 0
}
