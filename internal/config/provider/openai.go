package provider

import (
	"context"
	"fmt"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIProvider struct {
	client     *openai.Client
	embedModel string
	chatModel  string
}

func NewOpenAIProvider(apiKey string, embedModel string, chatModel string) *OpenAIProvider {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{
		client:     &client,
		embedModel: embedModel,
		chatModel:  chatModel,
	}
}

func (o *OpenAIProvider) Embed(text string) ([]float32, error) {
	ctx := context.Background()
	res, err := o.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbedding3Small, // or use the appropriate constant
		Input: openai.EmbeddingNewParamsInputArrayOfStrings([]string{text}),
	})
	if err != nil {
		return nil, err
	}
	if len(res.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	vec64 := res.Data[0].Embedding // []float64
	out := make([]float32, len(vec64))
	for i, v := range vec64 {
		out[i] = float32(v)
	}
	return out, nil
}
