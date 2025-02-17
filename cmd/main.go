package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

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

	resp, err := client.SendMessage("Hello, World!")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(resp.Content) > 0 {
		fmt.Printf("Response: %s\n", resp.Content[0].Text)
	}

	diffRange, err := git.NewDiffRange()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	changedFiles, err := diffRange.GetChangedFiles()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print(changedFiles)

}
