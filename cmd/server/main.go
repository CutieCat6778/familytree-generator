package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/server"
)

func main() {
	
	port := flag.Int("port", 8080, "Port to listen on")
	dataDir := flag.String("data", "./data", "Path to data directory")
	webDir := flag.String("web", "./web/dist", "Path to web static files directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Family Tree Generator API Server\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nAPI Endpoints:\n")
		fmt.Fprintf(os.Stderr, "  GET  /api/health          - Health check\n")
		fmt.Fprintf(os.Stderr, "  GET  /api/countries       - List available countries\n")
		fmt.Fprintf(os.Stderr, "  GET  /api/country/{slug}  - Get country statistics\n")
		fmt.Fprintf(os.Stderr, "  POST /api/generate        - Generate a family tree\n")
		fmt.Fprintf(os.Stderr, "\nGenerate Request Body (JSON):\n")
		fmt.Fprintf(os.Stderr, "  {\n")
		fmt.Fprintf(os.Stderr, "    \"country\": \"germany\",\n")
		fmt.Fprintf(os.Stderr, "    \"generations\": 3,\n")
		fmt.Fprintf(os.Stderr, "    \"seed\": 12345,\n")
		fmt.Fprintf(os.Stderr, "    \"start_year\": 1970,\n")
		fmt.Fprintf(os.Stderr, "    \"gender\": \"M\" or \"F\",\n")
		fmt.Fprintf(os.Stderr, "    \"include_extended\": false\n")
		fmt.Fprintf(os.Stderr, "  }\n")
	}

	flag.Parse()

	
	log.Printf("Loading data from %s...", *dataDir)
	repo, err := data.NewRepository(*dataDir)
	if err != nil {
		log.Fatalf("Error loading data: %v", err)
	}
	log.Printf("Data loaded successfully")

	
	if _, err := os.Stat(*webDir); os.IsNotExist(err) {
		log.Printf("Warning: Web directory %s not found. Run 'cd web && npm run build' to build the frontend.", *webDir)
		*webDir = "" 
	}

	
	addr := fmt.Sprintf(":%d", *port)
	srv := server.NewServer(repo, addr, *webDir)

	log.Printf("Starting server at http://localhost%s", addr)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
