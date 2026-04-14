package utils

import (
	"context"
	"fmt"

	"portfolio-website-backend/internal/config"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// LLMFactory provides instances of AI models
type LLMFactory struct{}

// NewLLMFactory creates a new LLMFactory
func NewLLMFactory() *LLMFactory {
	return &LLMFactory{}
}

// CreateGeminiClient initializes and returns the official Gemini client
func (f *LLMFactory) CreateGeminiClient(ctx context.Context) (*genai.Client, error) {
	if config.Envs.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	return genai.NewClient(ctx, option.WithAPIKey(config.Envs.GeminiAPIKey))
}

// EmbedQuery is a helper for single query embeddings
func (f *LLMFactory) EmbedQuery(ctx context.Context, text string, dimensions int) ([]float32, error) {
	client, err := f.CreateGeminiClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	modelName := config.Envs.GeminiEmbeddingModel
	if modelName == "" {
		modelName = "text-embedding-001" // default
	}

	em := client.EmbeddingModel(modelName)
	// Setting TaskType to RetrievalDocument as per previous python code
	em.TaskType = genai.TaskTypeRetrievalDocument
	// Although the SDK might not expose output dimensionality directly for all models in this version,
	// the standard vector size for this model is 768.

	res, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, err
	}

	if res != nil && res.Embedding != nil && len(res.Embedding.Values) > 0 {
		values := res.Embedding.Values
		if dimensions > 0 && len(values) > dimensions {
			values = values[:dimensions]
		}
		return values, nil
	}

	return nil, fmt.Errorf("failed to generate embeddings: no values returned")
}
