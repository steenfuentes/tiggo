package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/steenfuentes/tiggo/internal/analyze"
	"github.com/steenfuentes/tiggo/internal/cli"
	"github.com/steenfuentes/tiggo/internal/llm"
)

func main() {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	// Try loading .env file if API key not found in OS env
	if apiKey == "" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: No .env file found, checking system environment")
		}

		apiKey = os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			log.Fatal("NO API KEY FOUND in system environment or .env file")
		}
	}

	config, err := cli.ParseFlags()
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	diffRange, err := cli.BuildDiffRange(config)
	if err != nil {
		fmt.Printf("Error creating diff range: %v\n", err)
		os.Exit(1)
	}

	client := llm.NewLLMClient("anthropic", apiKey)

	analyzer := &analyze.Analyzer{
		LLMClient: client,
		DiffRange: diffRange,
	}

	if err := analyzer.DoAnalysis(); err != nil {
		fmt.Printf("Error performing analysis: %v\n", err)
		os.Exit(1)
	}

	if err := analyzer.Summarize("summary.md"); err != nil {
		fmt.Printf("Error writing summary: %v\n", err)
		os.Exit(1)
	}

}
