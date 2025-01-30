package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/prompts"
)

func PingOllama() error {
	host := viper.GetString("ollama.url")
	resp, err := http.Get(host + "/api/version")
	if err != nil {
		return fmt.Errorf("failed to ping ollama: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("ollama returned status code %d", resp.StatusCode)
	}

	return nil
}

type OllamaResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []int64   `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int64     `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int64     `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

func GenerateInsights(source string) (string, error) {
	prompt := prompts.NewPromptTemplate(
		"Summerize this post on Solar Network below: {{.content}}",
		[]string{"content"},
	)
	inPrompt, err := prompt.Format(map[string]any{
		"content": source,
	})
	if err != nil {
		return "", fmt.Errorf("failed to format prompt: %v", err)
	}

	raw, _ := json.Marshal(map[string]any{
		"model":  viper.GetString("ollama.model"),
		"prompt": inPrompt,
		"stream": false,
	})

	start := time.Now()

	url := viper.GetString("ollama.url") + "/api/generate"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(raw))
	if err != nil {
		return "", fmt.Errorf("failed to generate insights: %v", err)
	}
	outRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var response OllamaResponse
	err = json.Unmarshal(outRaw, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	took := time.Since(start)

	log.Info().Dur("took", took).Msg("Insight generated successfully...")

	return response.Response, err
}
