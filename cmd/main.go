package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/steenfuentes/gogit/internal/analyze"
	"github.com/steenfuentes/gogit/internal/cli"
	"github.com/steenfuentes/gogit/internal/llm"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("NO API KEY FOUND")
	}

	config, err := cli.ParseFlags()
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	// Build DiffRange from config
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
