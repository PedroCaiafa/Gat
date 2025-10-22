package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const openAIAPIBaseURL = "https://api.openai.com/v1"

type OpenAIProvider struct {
	apiKey     string
	embedModel string
	chatModel  string
	httpClient *http.Client
}

func NewOpenAIProvider(apiKey string, embedModel string, chatModel string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:     apiKey,
		embedModel: embedModel,
		chatModel:  chatModel,
		httpClient: &http.Client{},
	}
}

// ///////////// COISAS DE EMBEDDING  //////////////////////
type embeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type openAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func (o *OpenAIProvider) Embed(text string) ([]float32, error) {
	// Prepare the request body
	reqBody := embeddingRequest{
		Model: o.embedModel,
		Input: []string{text},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", openAIAPIBaseURL+"/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	// Send the request
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		var errResp openAIErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("OpenAI API error: %s (type: %s)", errResp.Error.Message, errResp.Error.Type)
	}

	// Parse the response
	var embResp embeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(embResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	// Convert []float64 to []float32
	vec64 := embResp.Data[0].Embedding
	out := make([]float32, len(vec64))
	for i, v := range vec64 {
		out[i] = float32(v)
	}

	return out, nil
}

///////////////////////////////////////////////////
