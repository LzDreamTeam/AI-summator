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

// Tests for the Add Tool
func TestAddTool(t *testing.T) {
	ctx := context.Background()

	input := &AddToolInput{A: 5.0, B: 3.0}
	result, err := AddTool(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "8.000000", result)
}

func TestAddTool_Float(t *testing.T) {
	ctx := context.Background()

	input := &AddToolInput{A: 1.5, B: 2.5}
	result, err := AddTool(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "4.000000", result)
}

func TestAddTool_NegativeNumbers(t *testing.T) {
	ctx := context.Background()

	input := &AddToolInput{A: -5.0, B: 3.0}
	result, err := AddTool(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "-2.000000", result)
}

func TestAddTool_Zero(t *testing.T) {
	ctx := context.Background()

	input := &AddToolInput{A: 0, B: 0}
	result, err := AddTool(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "0.000000", result)
}

func TestAddTool_NilInput(t *testing.T) {
	ctx := context.Background()

	result, err := AddTool(ctx, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "input cannot be nil")
}

func TestAddToolHandler_ValidJSON(t *testing.T) {
	ctx := context.Background()

	input := `{"a": 10.0, "b": 20.0}`
	result, err := AddToolHandler(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "30.000000", result)
}

func TestAddToolHandler_InvalidJSON(t *testing.T) {
	ctx := context.Background()

	input := `{"a": "not a number", "b": 20.0}`
	_, err := AddToolHandler(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse input")
}

func TestAddToolHandler_MalformedJSON(t *testing.T) {
	ctx := context.Background()

	input := `{invalid json}`
	_, err := AddToolHandler(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse input")
}
