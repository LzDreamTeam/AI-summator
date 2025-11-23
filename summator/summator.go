package summator

import (
	"context"
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
