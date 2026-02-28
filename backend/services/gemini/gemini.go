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
	if model == "" {
		return nil, fmt.Errorf("gemini model is required")
	}
	return &Client{
		client: client,
		model:  model,
	}, nil
}

// Categorize sends a list of food items (as a JSON string) to Gemini
// and returns the categorized JSON response.
func (c *Client) Categorize(ctx context.Context, foodItemsJSON string) (string, error) {
	prompt := CategorizerPrompt + "\n" + foodItemsJSON

	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}

	resp, err := c.client.Models.GenerateContent(ctx, c.model, contents, nil)
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

	resp, err := c.client.Models.GenerateContent(ctx, c.model, contents, nil)
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
