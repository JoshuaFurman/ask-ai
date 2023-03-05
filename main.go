package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func build_prompt(question string) string {
	return fmt.Sprintf(`
You are GPT-3 a highly advanced AI assistant that is an expert in all fields and you provide high quality answers.

I need your help answering the following question and please provides examples when you can:

%s
`, question)
}

func main() {
	// #########################################################################################
	// SET UP
	// #########################################################################################
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
	// #########################################################################################
	// #########################################################################################
	// #########################################################################################

	c := gogpt.NewClient(bearer)
	ctx := context.Background()

	// Make call to OpenAI API
	req := gogpt.CompletionRequest{
		Model:       "text-davinci-003",
		MaxTokens:   maxTokens,
		Temperature: float32(temperature),
		Prompt:      build_prompt(prompt), // Using build_prompt function to inject prompt template
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print results
	fmt.Println(resp.Choices[0].Text)
	fmt.Printf("\nTotal Tokens used: %d\n", resp.Usage.TotalTokens)
}
