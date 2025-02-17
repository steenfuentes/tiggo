package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/steenfuentes/gogit/internal/analyze"
	"github.com/steenfuentes/gogit/internal/git"
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

	client := llm.NewLLMClient("anthropic", apiKey)

	diffRange, err := git.NewDiffRange()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	analyzer := &analyze.Analyzer{
		LLMClient: client,
		DiffRange: diffRange,
	}

	err = analyzer.DoAnalysis()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = analyzer.Summarize("summary.md")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

}
