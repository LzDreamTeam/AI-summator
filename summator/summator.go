package summator

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

// LLMProvider defines the interface for generating content
type LLMProvider interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
}

// AISummator handles the summation using an LLM
type AISummator struct {
	llm LLMProvider
}

// NewAISummator creates a new instance of AISummator
func NewAISummator(llm LLMProvider) *AISummator {
	return &AISummator{
		llm: llm,
	}
}

// Description returns a description of the AISummator
func (s *AISummator) Description() string {
	return "Useful for adding two numbers together. Input should be a JSON object with 'a' and 'b' fields representing the numbers to add."
}

// Sum calculates the sum of two numbers using the LLM
func (s *AISummator) Sum(ctx context.Context, a, b float64) (float64, error) {
	prompt := fmt.Sprintf("Calculate the sum of %f and %f. Return ONLY the numeric result, nothing else. Do not include any text, just the number.", a, b)

	response, err := s.llm.Call(ctx, prompt)
	if err != nil {
		return 0, fmt.Errorf("failed to call LLM: %w", err)
	}

	// Clean up the response just in case the model adds whitespace or newlines
	cleanedResponse := strings.TrimSpace(response)

	result, err := strconv.ParseFloat(cleanedResponse, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse LLM response '%s' as float: %w", cleanedResponse, err)
	}

	return result, nil
}

// AddToolInput represents the input for the Add tool
type AddToolInput struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// AddTool is a simple tool that calculates a + b
func AddTool(ctx context.Context, input *AddToolInput) (string, error) {
	if input == nil {
		return "", fmt.Errorf("input cannot be nil")
	}
	result := input.A + input.B
	return fmt.Sprintf("%f", result), nil
}

// AddToolHandler parses JSON input and calls AddTool
func AddToolHandler(ctx context.Context, inputStr string) (string, error) {
	var input AddToolInput
	err := json.Unmarshal([]byte(inputStr), &input)
	if err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}
	return AddTool(ctx, &input)
}
