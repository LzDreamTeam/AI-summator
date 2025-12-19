package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"ai-summator/summator"
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
	// (Not used directly - tool is called instead of LLM)

	// Create a tool handler for the Add function
	toolInputA := summator.AddToolInput{A: a, B: b}
	
	// Call the Add tool directly instead of using LLM
	result, err := summator.AddTool(ctx, &toolInputA)
	if err != nil {
		log.Fatalf("Add tool failed: %v", err)
	}

	fmt.Printf("Result: %s\n", result)
}
