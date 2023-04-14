# OpenAI GPT-4 CLI Prompt
You must create an API account on OpenAI to receive a token: https://openai.com/

Currently the model used for completion is hardcoded as GPT-4 and as such you must have GPT-4 API access to use this program. This will be updated soon.

## Set Environment Variables:
- `OPEN_AI_TOKEN`: Your OpenAI API key (string)
- `OPEN_AI_TEMP`: Temperature value for API call (string)
- `OPEN_AI_MAX_TOKENS`: Max token value for API call (string)

## To Build Executable:
Run the following command from this directory:

`go build -o <executable-name>` and then add your executable to your PATH

## To Run:
- One-off request: `<executable-name> 'Explain Big O notation.'`

- Chat Mode: `<executable-name> --chat` or `<executable-name> -c`