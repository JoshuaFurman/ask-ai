# OpenAI GPT-4 CLI Prompt

You must create an API account on OpenAI to receive a token: https://openai.com/

Currently the model used for completion is hardcoded as GPT-4 and as such you must have GPT-4 API access to use this program. This will be updated soon.

## Set Environment Variables:

- `OPEN_AI_TOKEN`: Your OpenAI API key (string)
- `OPEN_AI_TEMP`: Temperature value for API call (string)
- `OPEN_AI_MAX_TOKENS`: Max token value for API call (string)
- `OPEN_AI_MODEL`: OpenAI Chat completions model to use (string)
- `OPEN_AI_CUSTOM_URL`: Custom URL if server your own model

## Acceptable Models

- "gpt-4-32k-0314"
- "gpt-4-32k"
- "gpt-4-0314"
- "gpt-4"
- "gpt-3.5-turbo-0301"
- "gpt-3.5-turbo"
- "mistralai/Mixtral-8x7B-Instruct-v0.1"

## To Build Executable:
Run the following command from this directory:

`go build -o <executable-name>` and then add your executable to your PATH

Please use the `update` script if you would like to pull code changes and build the ask binary. 

## To Run:
- One-off request: `<executable-name> 'Explain Big O notation.'`

- Chat Mode: `<executable-name> --chat` or `<executable-name> -c`

## Pasting Code in Chat Mode:
Something to note about pasting code into the terminal when using chat mode... Code **MUST BE** wrapped in triple-ticks: ``` otherwise the reader will stop accepting input and send the request to the OpenAI API without your full message. This is the cleanest solution I could come up with for reading code pasted into the terminal. 