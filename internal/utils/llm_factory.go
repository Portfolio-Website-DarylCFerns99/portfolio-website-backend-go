package utils

// EmbeddingsModel represents a generator for text embeddings
type EmbeddingsModel interface {
	EmbedQuery(text string, dimensions int) ([]float32, error)
}

// LLMFactory provides instances of AI models
type LLMFactory struct{}

// CreateEmbeddingsModel returns the appropriate embeddings model based on the provider string
func (f *LLMFactory) CreateEmbeddingsModel(provider string) EmbeddingsModel {
	// In a complete implementation, this would return a Gemini API or OpenAI client wrapper.
	// For now, it returns a stub or un-implemented version.
	return &StubEmbeddingsModel{}
}

// StubEmbeddingsModel is a temporary stub for vector generation
type StubEmbeddingsModel struct{}

// EmbedQuery stub generates an empty or dummy vector
func (m *StubEmbeddingsModel) EmbedQuery(text string, dimensions int) ([]float32, error) {
	vec := make([]float32, dimensions)
	// Example: fill with tiny non-zero values or leave zeroes (will not work for meaningful search)
	return vec, nil
}
