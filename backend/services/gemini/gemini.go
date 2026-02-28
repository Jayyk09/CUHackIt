package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// Client wraps the Google GenAI client for Gemini API interactions.
type Client struct {
	client *genai.Client
	model  string
	models []string
}

// New creates a new Gemini client using the provided API key and model.
// If model is empty, a default Gemini model is selected.
func New(ctx context.Context, apiKey string, model string) (*Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	models, err := discoverModels(ctx, client, model)
	if model == "" && len(models) > 0 {
		model = models[0]
	}
	return &Client{
		client: client,
		model:  model,
		models: models,
	}, nil
}

// Categorize sends a list of food items (as a JSON string) to Gemini
// and returns the categorized JSON response.
func (c *Client) Categorize(ctx context.Context, foodItemsJSON string) (string, error) {
	prompt := CategorizerPrompt + "\n" + foodItemsJSON

	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}

	resp, err := c.generateContentWithFallback(ctx, contents)
	if err != nil {
		return "", fmt.Errorf("gemini categorize error: %w", err)
	}

	return extractText(resp), nil
}

// Recommend sends a pantry list (as a JSON string) to Gemini
// and returns recipe recommendations as a JSON response.
func (c *Client) Recommend(ctx context.Context, pantryJSON string) (string, error) {
	prompt := RecipePrompt + "\n" + pantryJSON

	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}

	resp, err := c.generateContentWithFallback(ctx, contents)
	if err != nil {
		return "", fmt.Errorf("gemini recommend error: %w", err)
	}

	return extractText(resp), nil
}

// extractText pulls the text content from a GenerateContentResponse.
func extractText(resp *genai.GenerateContentResponse) string {
	if resp == nil || len(resp.Candidates) == 0 {
		return ""
	}
	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return ""
	}
	var result string
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			result += part.Text
		}
	}
	return result
}

func (c *Client) generateContentWithFallback(ctx context.Context, contents []*genai.Content) (*genai.GenerateContentResponse, error) {
	var lastErr error
	for _, model := range c.models {
		resp, err := c.client.Models.GenerateContent(ctx, model, contents, nil)
		if err == nil {
			c.model = model
			return resp, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func discoverModels(ctx context.Context, client *genai.Client, preferred string) ([]string, error) {
	models := make([]string, 0)
	for item, err := range client.Models.All(ctx) {
		if err != nil {
			return nil, err
		}
		if item == nil || item.Name == "" {
			continue
		}
		if !supportsGenerateContent(item.SupportedActions) {
			continue
		}
		models = append(models, item.Name)
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("no gemini models available")
	}
	if preferred == "" {
		return models, nil
	}
	ordered := []string{preferred}
	for _, name := range models {
		if name != preferred {
			ordered = append(ordered, name)
		}
	}
	return ordered, nil
}

func supportsGenerateContent(actions []string) bool {
	for _, action := range actions {
		if action == "generateContent" {
			return true
		}
	}
	return false
}
