package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func main() {
	bearer := os.Getenv("OPEN_AI_TOKEN")

	preTemp := os.Getenv("OPEN_AI_TEMP") // Pull env var, convert and error check
	temperature, err := strconv.ParseFloat(preTemp, 32)
	if err != nil {
		fmt.Println(err)
	}

	preMaxTokens := os.Getenv("OPEN_AI_MAX_TOKENS") // Pull env var, convert and error check
	maxTokens, err := strconv.Atoi(preMaxTokens)
	if err != nil {
		fmt.Println(err)
	}

	prompt := os.Args[1] // CLI argument provided as prompt to GPT-3

	c := gogpt.NewClient(bearer)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:       "text-davinci-003",
		MaxTokens:   maxTokens,
		Temperature: float32(temperature),
		Prompt:      prompt,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return
	}
	fmt.Println(resp.Choices[0].Text)
}
