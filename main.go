package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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

	if len(os.Args) < 2 {
		fmt.Println()
		fmt.Println("No command line argument supplied!")
		fmt.Println("Please use: --chat or -c for chat mode. Type \"exit\" to exit the chat.")
		fmt.Println("Please use: \"Double Quotes\" or 'Single Quotes' for a one-off command.")
		return

	} else if os.Args[1] == "--chat" || os.Args[1] == "-c" {
		client := openai.NewClient(bearer)
		ctx := context.Background()
		messages := make([]openai.ChatCompletionMessage, 0)
		reader := bufio.NewReader(os.Stdin)
		fmt.Println()
		fmt.Println("Begin Conversation")
		fmt.Println("------------------")
		var fullMessage string

		for {
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)
			if text == "exit" {
				return
			}
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			})

			stream, err := client.CreateChatCompletionStream(
				ctx,
				openai.ChatCompletionRequest{
					Model:       openai.GPT4,
					MaxTokens:   maxTokens,
					Temperature: float32(temperature),
					Messages:    messages,
					Stream:      true,
				},
			)

			if err != nil {
				fmt.Printf("ChatCompletion error: %v\n", err)
				continue
			}
			defer stream.Close()

			// Used for processing code blocks
			state := 0
			buffer := ""
			code := ""
			lang := ""

			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					fmt.Println() // For spacing of the respsonse
					fmt.Println()
					break
				}
				if err != nil {
					fmt.Printf("\nStream error: %v\n", err)
					return
				}

				for _, char := range response.Choices[0].Delta.Content {
					processChar(char, &state, &buffer, &code, &lang)
				}
				fullMessage = fullMessage + response.Choices[0].Delta.Content // Save full response to save back to chat context

			}
			// Put response back into chat context
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: fullMessage,
			})
		}

	} else if os.Args[1][0] == '"' || os.Args[1][0] == '\'' { // Checks the first char of the cli argument
		// Get prompt and validate
		system_prompt := "You are a super powerful AI assistant. Answer all queries as concisely as possible and try to think through each response step-by-step."
		prompt := os.Args[1]

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
				fmt.Println() // For spacing of the response
				return        // Finished displaying stream to terminal
			}
			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return
			}

			for _, char := range response.Choices[0].Delta.Content {
				processChar(char, &state, &buffer, &code, &lang)
			}
		}
	} else {
		fmt.Println()
		fmt.Println("***ERROR***")
		fmt.Println("Sorry that was not a correct argument.")
		fmt.Println("Please use: --chat or -c for chat mode. Type \"exit\" to exit the chat.")
		fmt.Println("Please use: \"Double Quotes\" or 'Single Quotes' for a one-off command.")
		return
	}
}
