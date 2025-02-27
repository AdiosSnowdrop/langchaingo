package ollama

import (
	"context"
	"strings"

	"github.com/adiossnowdrop/langchaingo/embeddings"
	"github.com/adiossnowdrop/langchaingo/llms/ollama"
)

// Ollama is the embedder using the Ollama api.
type Ollama struct {
	client *ollama.LLM

	StripNewLines bool
	BatchSize     int
}

var _ embeddings.Embedder = Ollama{}

// NewOllama creates a new Ollama with options. Options for client, strip new lines and batch.
func NewOllama(opts ...Option) (Ollama, error) {
	o, err := applyClientOptions(opts...)
	if err != nil {
		return Ollama{}, err
	}

	return o, nil
}

// EmbedDocuments creates one vector embedding for each of the texts.
func (e Ollama) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	batchedTexts := embeddings.BatchTexts(
		embeddings.MaybeRemoveNewLines(texts, e.StripNewLines),
		e.BatchSize,
	)

	emb := make([][]float32, 0, len(texts))
	for _, texts := range batchedTexts {
		curTextEmbeddings, err := e.client.CreateEmbedding(ctx, texts)
		if err != nil {
			return nil, err
		}

		textLengths := make([]int, 0, len(texts))
		for _, text := range texts {
			textLengths = append(textLengths, len(text))
		}

		combined, err := embeddings.CombineVectors(curTextEmbeddings, textLengths)
		if err != nil {
			return nil, err
		}

		emb = append(emb, combined)
	}

	return emb, nil
}

// EmbedQuery embeds a single text.
func (e Ollama) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	if e.StripNewLines {
		text = strings.ReplaceAll(text, "\n", " ")
	}

	emb, err := e.client.CreateEmbedding(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	return emb[0], nil
}
