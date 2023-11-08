package agents

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/adiossnowdrop/langchaingo/chains"
	"github.com/adiossnowdrop/langchaingo/llms/openai"
	"github.com/adiossnowdrop/langchaingo/memory"
	"github.com/adiossnowdrop/langchaingo/tools"
	"github.com/stretchr/testify/require"
)

func TestConversationalWithMemory(t *testing.T) {
	t.Parallel()
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	llm, err := openai.New(openai.WithModel("gpt-4"))
	require.NoError(t, err)

	executor, err := Initialize(
		llm,
		[]tools.Tool{tools.Calculator{}},
		ConversationalReactDescription,
		WithMemory(memory.NewConversationBuffer()),
	)
	require.NoError(t, err)

	_, err = chains.Run(context.Background(), executor, "Hi! my name is Bob and the year I was born is 1987")
	require.NoError(t, err)

	res, err := chains.Run(context.Background(), executor, "What is the year I was born times 34")
	require.NoError(t, err)
	require.True(t, strings.Contains(res, "67558"), `result does not contain the correct answer '67558'`)
}
