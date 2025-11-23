package summator

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms/ollama"
)

func TestIntegrationSum(t *testing.T) {
	// This test assumes Ollama is running locally with llama3.1 model pulled.
	// To run this test: go test -v ./summator/ -run TestIntegrationSum
	
	llm, err := ollama.New(ollama.WithModel("llama3.1"))
	if err != nil {
		t.Skipf("Skipping integration test: failed to create ollama client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	summator := NewAISummator(llm)
	
	// We use a simple sum. Note: LLMs are probabilistic, so 10+20 might not always be 30 in text,
	// but with the prompt we designed, it should be robust.
	result, err := summator.Sum(ctx, 10, 20)
	if err != nil {
		t.Logf("Integration test failed (ensure ollama is running and llama3.1 is pulled): %v", err)
		t.FailNow()
	}

	assert.Equal(t, 30.0, result)
}
