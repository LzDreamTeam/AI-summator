package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"ai-summator/summator"

	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ai-summator <number1> <number2>")
		fmt.Println("Example: ai-summator 5.5 3.2")
		os.Exit(1)
	}

	a, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		log.Fatalf("Invalid first number: %v", err)
	}

	b, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatalf("Invalid second number: %v", err)
	}

	ctx := context.Background()

	// Initialize Ollama with llama3.1
	// Ensure you have ollama running and llama3.1 pulled: `ollama pull llama3.1`
	llm, err := ollama.New(ollama.WithModel("llama3.1"))
	if err != nil {
		log.Fatalf("Failed to create Ollama client: %v", err)
	}

	// Create the Add tool using AISummator with its description
	app := summator.NewAISummator(llm)
	toolDescription := app.Description()

	// Create input for the Add tool
	toolInput := summator.AddToolInput{A: a, B: b}

	// Call the Add tool directly
	result, err := summator.AddTool(ctx, &toolInput)
	if err != nil {
		log.Fatalf("Add tool failed: %v", err)
	}

	fmt.Printf("Tool Name: Add\n")
	fmt.Printf("Tool Description: %s\n", toolDescription)
	fmt.Printf("Input: {\"a\": %f, \"b\": %f}\n", a, b)
	fmt.Printf("Result: %s\n", result)
}
