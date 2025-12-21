package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"ai-summator/summator"

	"github.com/tmc/langchaingo/llms"
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

	// Create agent executor with Add tool
	executor := createAgentExecutor(llm, a, b)

	// Create prompt for the agent
	userQuery := fmt.Sprintf("Calculate the sum of %f and %f", a, b)

	// Invoke the agent
	result, err := executor(ctx, userQuery)
	if err != nil {
		log.Fatalf("Agent failed: %v", err)
	}

	fmt.Printf("Result: %s\n", result)
}

// Tool represents a callable tool with description
type Tool struct {
	Name        string
	Description string
	Func        func(context.Context, string) (string, error)
}

// createAgentExecutor creates an agent that can use the Add tool
func createAgentExecutor(llm llms.Model, a, b float64) func(context.Context, string) (string, error) {
	// Define the Add tool
	addTool := Tool{
		Name:        "Add",
		Description: "Calculates the sum of two numbers. Takes 'a' and 'b' as input and returns their sum.",
		Func: func(ctx context.Context, input string) (string, error) {
			return summator.AddToolHandler(ctx, input)
		},
	}

	tools := map[string]Tool{"Add": addTool}

	return func(ctx context.Context, query string) (string, error) {
		// Call LLM to determine which tool to use
		prompt := fmt.Sprintf(`You are an agent that can use tools to answer questions.
You have access to the following tool:
- Add: %s

User query: %s

Based on the user query, decide which tool to use and provide your answer.
If you use a tool, respond in JSON format: {"tool": "tool_name", "input": {"a": value, "b": value}}
If you can answer without tools, just respond with the answer directly.

Your response:`, addTool.Description, query)

		response, err := llm.Call(ctx, prompt, llms.WithTemperature(0))
		if err != nil {
			return "", fmt.Errorf("failed to call LLM: %w", err)
		}

		// Try to parse JSON response to see if a tool was requested
		var toolCall map[string]interface{}
		err = json.Unmarshal([]byte(response), &toolCall)
		if err == nil && toolCall["tool"] != nil {
			toolName := toolCall["tool"].(string)
			if tool, ok := tools[toolName]; ok {
				// Convert input to JSON string
				inputJSON, err := json.Marshal(toolCall["input"])
				if err != nil {
					return "", fmt.Errorf("failed to marshal tool input: %w", err)
				}
				// Call the tool
				return tool.Func(ctx, string(inputJSON))
			}
		}

		// If no tool was called, return the LLM response
		return response, nil
	}
}
