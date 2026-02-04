package config

import (
	"github.com/familytree-generator/internal/generator"
	"github.com/familytree-generator/internal/model"
)

// AppConfig holds all application configuration
type AppConfig struct {
	// Generation settings
	Country         string
	Generations     int
	Seed            int64
	StartYear       int
	RootGender      string
	IncludeExtended bool

	// Output settings
	OutputPath   string
	OutputFormat string // "csv", "json", or "both"

	// Data settings
	DataDir string

	// Flags
	ListCountries bool
	Verbose       bool
}

// DefaultAppConfig returns the default application configuration
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		Country:         "germany",
		Generations:     3,
		Seed:            0, // 0 means random
		StartYear:       1970,
		RootGender:      "random",
		IncludeExtended: false,
		OutputPath:      "family_tree.csv",
		OutputFormat:    "csv",
		DataDir:         "./data",
		ListCountries:   false,
		Verbose:         false,
	}
}

// ToGeneratorConfig converts app config to generator config
func (c *AppConfig) ToGeneratorConfig() generator.Config {
	var gender model.Gender
	switch c.RootGender {
	case "M", "m", "male":
		gender = model.Male
	case "F", "f", "female":
		gender = model.Female
	default:
		gender = "" // Random
	}

	return generator.Config{
		Country:         c.Country,
		Generations:     c.Generations,
		Seed:            c.Seed,
		StartYear:       c.StartYear,
		RootGender:      gender,
		IncludeExtended: c.IncludeExtended,
	}
}

// Validate checks if the configuration is valid
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

	return nil
}
