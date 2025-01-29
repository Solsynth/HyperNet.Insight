package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
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

var LargeModel *ollama.LLM

func ConnectOllama() error {
	model := viper.GetString("ollama.model")
	llm, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		return err
	}
	LargeModel = llm
	return nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	completion, err := LargeModel.Call(ctx, inPrompt,
		llms.WithTemperature(0.8),
	)

	return completion, err
}
