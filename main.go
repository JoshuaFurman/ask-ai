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
	if bearer == "" {
		fmt.Println("You need to set OPEN_AI_TOKEN environment variable.")
	}

	preTemp := os.Getenv("OPEN_AI_TEMP") // Pull env var, convert and error check
	if preTemp == "" {
		fmt.Println("You need to set OPEN_AI_TEMP environment variable.")
	}
	temperature, err := strconv.ParseFloat(preTemp, 32)
	if err != nil {
		fmt.Println(err)
	}

	preMaxTokens := os.Getenv("OPEN_AI_MAX_TOKENS") // Pull env var, convert and error check
	if preMaxTokens == "" {
		fmt.Println("You need to set OPEN_AI_MAX_TOKENS evironment variable.")
	}
	maxTokens, err := strconv.Atoi(preMaxTokens)
	if err != nil {
		fmt.Println(err)
	}

	var prompt string
	if len(os.Args) < 2 {
		fmt.Println("No command line argument supplied!")
		return
	} else {
		prompt = os.Args[1] // CLI argument provided as prompt to GPT-3
	}

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
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Choices[0].Text)
}
