package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Configure the API client
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

	// Get prompt and validate
	var prompt string
	if len(os.Args) < 2 {
		fmt.Println("No command line argument supplied!")
		return
	} else {
		prompt = os.Args[1] // CLI argument provided as prompt to GPT-4
	}

	// Create GPT stream chat request
	c := openai.NewClient(bearer)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4,
		MaxTokens:   maxTokens,
		Temperature: float32(temperature),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Println()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// fmt.Println("\nStream finished")
			fmt.Println()
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
