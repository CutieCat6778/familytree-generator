package data

import (
	"fmt"
	"path/filepath"
)

type DemographicData struct {
	BirthRates           map[string]float64
	DeathRates           map[string]float64
	LifeExpectancy       map[string]float64
	LifeExpectancyFemale map[string]float64
	LifeExpectancyMale   map[string]float64
	MigrationRates       map[string]float64
	InfantMortality      map[string]float64
	Population           map[string]float64
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

	if err := d.loadLifeExpectancy(filepath.Join(dataDir, "life_exp_at_birth_by_sex.csv"), filepath.Join(dataDir, "life_exp_at_birth.csv")); err != nil {
		return nil, err
	}

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

type LifeExpectancyBySex struct {
	Total  float64
	Female float64
	Male   float64
}

func (d *DemographicData) loadLifeExpectancy(bySexPath, legacyPath string) error {
	bySex, err := LoadLifeExpectancyBySexCSV(bySexPath)
	if err == nil && len(bySex) > 0 {
		d.LifeExpectancy = make(map[string]float64, len(bySex))
		d.LifeExpectancyFemale = make(map[string]float64, len(bySex))
		d.LifeExpectancyMale = make(map[string]float64, len(bySex))
		for slug, v := range bySex {
			if v.Total > 0 {
				d.LifeExpectancy[slug] = v.Total
			}
			if v.Female > 0 {
				d.LifeExpectancyFemale[slug] = v.Female
			}
			if v.Male > 0 {
				d.LifeExpectancyMale[slug] = v.Male
			}
		}
		records, legacyErr := LoadCSV(legacyPath)
		if legacyErr == nil {
			legacy := RecordsToMap(records)
			for slug, value := range legacy {
				if _, ok := d.LifeExpectancy[slug]; !ok && value > 0 {
					d.LifeExpectancy[slug] = value
				}
			}
		}
		return nil
	}

	records, err := LoadCSV(legacyPath)
	if err != nil {
		return fmt.Errorf("loading life expectancy data: %w", err)
	}
	d.LifeExpectancy = RecordsToMap(records)
	return nil
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

func (d *DemographicData) GetLifeExpectancyFemale(slug string) float64 {
	if v, ok := d.LifeExpectancyFemale[slug]; ok {
		return v
	}
	return d.GetLifeExpectancy(slug)
}

func (d *DemographicData) GetLifeExpectancyMale(slug string) float64 {
	if v, ok := d.LifeExpectancyMale[slug]; ok {
		return v
	}
	return d.GetLifeExpectancy(slug)
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
