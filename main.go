package main

import (
	"context"
	"fmt"
	"os"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func main() {
	bearer := os.Getenv("OPEN_AI_TOKEN")

	prompt := os.Args[1]

	c := gogpt.NewClient(bearer)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:       "text-davinci-002",
		MaxTokens:   500,
		Temperature: 0.5,
		Prompt:      prompt,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return
	}
	fmt.Println(resp.Choices[0].Text)
}
