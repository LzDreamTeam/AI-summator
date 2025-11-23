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

	// Initialize Ollama with phi3
	// Ensure you have ollama running and phi3 pulled: `ollama pull phi3`
	llm, err := ollama.New(ollama.WithModel("phi3"))
	if err != nil {
		log.Fatalf("Failed to create Ollama client: %v", err)
	}

	app := summator.NewAISummator(llm)
	result, err := app.Sum(ctx, a, b)
	if err != nil {
		log.Fatalf("Summation failed: %v", err)
	}

	fmt.Printf("Result: %f\n", result)
}
