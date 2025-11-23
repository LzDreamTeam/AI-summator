package summator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/llms"
)

// MockLLM is a mock implementation of LLMProvider
type MockLLM struct {
	mock.Mock
}

func (m *MockLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	// We need to handle the variadic arguments for options
	// Since testify mock doesn't handle variadic args directly in the signature matching easily without specific setup,
	// we will just pass them through.
	// For simplicity in this test, we ignore the options in the Called signature or pass them as a slice.
	// However, standard Called usage:
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func TestSum(t *testing.T) {
	mockLLM := new(MockLLM)
	summator := NewAISummator(mockLLM)
	ctx := context.Background()

	// Expectation
	expectedPrompt := "Calculate the sum of 5.000000 and 3.000000. Return ONLY the numeric result, nothing else. Do not include any text, just the number."
	mockLLM.On("Call", ctx, expectedPrompt).Return("8", nil)

	result, err := summator.Sum(ctx, 5, 3)

	assert.NoError(t, err)
	assert.Equal(t, 8.0, result)
	mockLLM.AssertExpectations(t)
}

func TestSum_Float(t *testing.T) {
	mockLLM := new(MockLLM)
	summator := NewAISummator(mockLLM)
	ctx := context.Background()

	// Expectation
	expectedPrompt := "Calculate the sum of 1.500000 and 2.500000. Return ONLY the numeric result, nothing else. Do not include any text, just the number."
	mockLLM.On("Call", ctx, expectedPrompt).Return("4.0", nil)

	result, err := summator.Sum(ctx, 1.5, 2.5)

	assert.NoError(t, err)
	assert.Equal(t, 4.0, result)
	mockLLM.AssertExpectations(t)
}

func TestSum_Error(t *testing.T) {
	mockLLM := new(MockLLM)
	summator := NewAISummator(mockLLM)
	ctx := context.Background()

	mockLLM.On("Call", ctx, mock.Anything).Return("not a number", nil)

	_, err := summator.Sum(ctx, 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse LLM response")
}
