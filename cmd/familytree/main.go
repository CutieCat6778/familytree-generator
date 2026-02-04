package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/familytree-generator/internal/config"
	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/generator"
	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/internal/output"
)

func main() {
	cfg := config.DefaultAppConfig()

	
	flag.StringVar(&cfg.Country, "country", cfg.Country, "Country slug for demographics (e.g., 'united-states', 'japan')")
	flag.IntVar(&cfg.Generations, "generations", cfg.Generations, "Number of generations to generate (1-10)")
	flag.Int64Var(&cfg.Seed, "seed", cfg.Seed, "Random seed for reproducibility (0 = random)")
	flag.StringVar(&cfg.OutputPath, "output", cfg.OutputPath, "Output file path")
	flag.StringVar(&cfg.OutputFormat, "format", cfg.OutputFormat, "Output format: csv, json, or both")
	flag.StringVar(&cfg.DataDir, "data", cfg.DataDir, "Path to data directory")
	flag.BoolVar(&cfg.ListCountries, "list-countries", cfg.ListCountries, "List available countries and exit")
	flag.IntVar(&cfg.StartYear, "start-year", cfg.StartYear, "Birth year of the root person")
	flag.StringVar(&cfg.RootGender, "gender", cfg.RootGender, "Root person gender: M, F, or random")
	flag.BoolVar(&cfg.IncludeExtended, "extended", cfg.IncludeExtended, "Include extended family (siblings)")
	flag.BoolVar(&cfg.Verbose, "verbose", cfg.Verbose, "Enable verbose output")

	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Family Tree Generator - Generate realistic family trees based on demographic data\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -country japan -generations 5 -seed 12345\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -country germany -format json -output tree.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -list-countries\n", os.Args[0])
	}

	flag.Parse()

	
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	
	if cfg.Verbose {
		fmt.Printf("Loading data from %s...\n", cfg.DataDir)
	}

	repo, err := data.NewRepository(cfg.DataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading data: %v\n", err)
		os.Exit(1)
	}

	
	if cfg.ListCountries {
		listCountries(repo)
		return
	}

	
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}

	
	if err := repo.ValidateCountry(cfg.Country); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Use -list-countries to see available countries\n")
		os.Exit(1)
	}

	
	genConfig := cfg.ToGeneratorConfig()
	engine := generator.NewEngine(genConfig, repo)

	if cfg.Verbose {
		fmt.Printf("Generating family tree for %s...\n", cfg.Country)
		fmt.Printf("  Generations: %d\n", cfg.Generations)
		fmt.Printf("  Start year: %d\n", cfg.StartYear)
		fmt.Printf("  Seed: %d\n", cfg.Seed)
		fmt.Printf("  Extended family: %v\n", cfg.IncludeExtended)
	}

	
	tree, err := engine.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating tree: %v\n", err)
		os.Exit(1)
	}

	if cfg.Verbose {
		fmt.Printf("Generated %d persons in %d families\n", tree.PersonCount(), tree.FamilyCount())
	}

	
	if err := writeOutput(tree, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Family tree generated successfully!\n")
	fmt.Printf("  Persons: %d\n", tree.PersonCount())
	fmt.Printf("  Families: %d\n", tree.FamilyCount())
	fmt.Printf("  Seed: %d (use this to reproduce the same tree)\n", cfg.Seed)
}

func listCountries(repo *data.Repository) {
	countries := repo.GetCountriesWithNames()

	fmt.Printf("Available countries with complete data (%d):\n\n", len(countries))

	
	sort.Strings(countries)

	for i, slug := range countries {
		fmt.Printf("  %-35s", slug)
		if (i+1)%2 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()

	fmt.Printf("\nNote: Use the slug (lowercase with dashes) with the -country flag.\n")
	fmt.Printf("Example: familytree -country united-states\n")
}

func writeOutput(tree *model.FamilyTree, cfg *config.AppConfig) error {
	format := strings.ToLower(cfg.OutputFormat)

	switch format {
	case "csv":
		if err := output.WriteCSV(tree, cfg.OutputPath); err != nil {
			return err
		}
		fmt.Printf("Output written to: %s\n", cfg.OutputPath)

	case "json":
		if err := output.WriteJSON(tree, cfg.OutputPath); err != nil {
			return err
		}
		fmt.Printf("Output written to: %s\n", cfg.OutputPath)

		
		vizPath := strings.TrimSuffix(cfg.OutputPath, filepath.Ext(cfg.OutputPath)) + "_viz.json"
		if err := output.WriteVisualizationJSON(tree, vizPath); err != nil {
			return err
		}
		fmt.Printf("Visualization data written to: %s\n", vizPath)

	case "both":
		
		csvPath := strings.TrimSuffix(cfg.OutputPath, filepath.Ext(cfg.OutputPath)) + ".csv"
		if err := output.WriteCSV(tree, csvPath); err != nil {
			return err
		}
		fmt.Printf("CSV output written to: %s\n", csvPath)

		
		jsonPath := strings.TrimSuffix(cfg.OutputPath, filepath.Ext(cfg.OutputPath)) + ".json"
		if err := output.WriteJSON(tree, jsonPath); err != nil {
			return err
		}
		fmt.Printf("JSON output written to: %s\n", jsonPath)

		
		vizPath := strings.TrimSuffix(cfg.OutputPath, filepath.Ext(cfg.OutputPath)) + "_viz.json"
		if err := output.WriteVisualizationJSON(tree, vizPath); err != nil {
			return err
		}
		fmt.Printf("Visualization data written to: %s\n", vizPath)

	default:
		return fmt.Errorf("unknown output format: %s", format)
	}

	return nil
}
