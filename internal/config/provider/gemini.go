package provider

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	apiKey     string
	embedModel string
	chatModel  string
}

func NewGeminiProvider(apiKey string, embedModel string, chatModel string) *GeminiProvider {
	return &GeminiProvider{
		apiKey:     apiKey,
		embedModel: embedModel,
		chatModel:  chatModel,
	}
}

func (g *GeminiProvider) Embed(text string) ([]float32, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  g.apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	contents := []*genai.Content{genai.NewContentFromText(text, genai.RoleUser)}

	result, err := client.Models.EmbedContent(ctx,
		g.embedModel,
		contents,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	if len(result.Embeddings) == 0 || len(result.Embeddings[0].Values) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	// Convert to []float32
	embedding := result.Embeddings[0].Values
	out := make([]float32, len(embedding))
	for i, v := range embedding {
		out[i] = float32(v)
	}

	return out, nil
}
