package config

import (
	"github.com/familytree-generator/internal/generator"
	"github.com/familytree-generator/internal/model"
)

type AppConfig struct {
	Country            string
	Generations        int
	Seed               int64
	StartYear          int
	RootGender         string
	IncludeExtended    bool
	LifeExpectancyMode string

	OutputPath   string
	OutputFormat string

	DataDir string

	ListCountries bool
	Verbose       bool
}

func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		Country:            "germany",
		Generations:        3,
		Seed:               0,
		StartYear:          1970,
		RootGender:         "random",
		IncludeExtended:    false,
		LifeExpectancyMode: string(generator.LifeExpectancyTotal),
		OutputPath:         "family_tree.csv",
		OutputFormat:       "csv",
		DataDir:            "./data",
		ListCountries:      false,
		Verbose:            false,
	}
}

func (c *AppConfig) ToGeneratorConfig() generator.Config {
	var gender model.Gender
	switch c.RootGender {
	case "M", "m", "male":
		gender = model.Male
	case "F", "f", "female":
		gender = model.Female
	default:
		gender = ""
	}

	return generator.Config{
		Country:            c.Country,
		Generations:        c.Generations,
		Seed:               c.Seed,
		StartYear:          c.StartYear,
		RootGender:         gender,
		IncludeExtended:    c.IncludeExtended,
		LifeExpectancyMode: generator.ParseLifeExpectancyMode(c.LifeExpectancyMode),
	}
}

func (c *AppConfig) Validate() error {
	if c.Generations < 1 {
		c.Generations = 1
	}
	if c.Generations > 10 {
		c.Generations = 10
	}

	if c.StartYear < 1800 {
		c.StartYear = 1800
	}
	if c.StartYear > 2024 {
		c.StartYear = 2024
	}

	c.LifeExpectancyMode = string(generator.ParseLifeExpectancyMode(c.LifeExpectancyMode))

	return nil
}
