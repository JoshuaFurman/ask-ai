package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	openai "github.com/sashabaranov/go-openai"
)

func formatCode(code, lang string) {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("solarized-dark")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	err = formatter.Format(os.Stdout, style, iterator)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func processChar(char rune, state *int, buffer *string, code *string, lang *string) {
	switch *state {
	case 0:
		if char == '`' {
			*state = 1
		} else {
			fmt.Print(string(char))
		}
	case 1:
		if char == '`' {
			*state = 2
		} else {
			*state = 0
			fmt.Print("`", string(char))
		}
	case 2:
		if char == '`' {
			*state = 3
			*buffer = ""
		} else {
			*state = 0
			fmt.Print("``", string(char))
		}
	case 3:
		if char == '\n' {
			*state = 4
			*lang = *buffer
			*buffer = ""
		} else {
			*buffer += string(char)
		}
	case 4:
		if char == '`' {
			*state = 5
		} else {
			*code += string(char)
		}
	case 5:
		if char == '`' {
			*state = 6
		} else {
			*state = 4
			*code += "`" + string(char)
		}
	case 6:
		if char == '`' {
			formatCode(*code, *lang)
			*code = ""
			*lang = ""
			*state = 0
		} else {
			*state = 4
			*code += "``" + string(char)
		}
	}
}

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
	system_prompt := "You are a super powerful AI assistant. Answer all queries as concisely as possible and try to think through each response step-by-step."
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
				Role:    openai.ChatMessageRoleSystem,
				Content: system_prompt,
			},
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

	state := 0
	buffer := ""
	code := ""
	lang := ""

	fmt.Println()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println() // For spacing of the respsonse
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		for _, char := range response.Choices[0].Delta.Content {
			processChar(char, &state, &buffer, &code, &lang)
		}
	}
}
