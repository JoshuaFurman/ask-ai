# OpenAI GPT-3 CLI Prompt
You must create an API account on OpenAI to receive a token: https://openai.com/

## Set Environment Variables:
- ```OPEN_AI_TOKEN```: Your OpenAI API key
- ```OPEN_AI_TEMP```: Temperature value for API call
- ```OPEN_AI_MAX_TOKENS```: Max token value for API call

## To Build Executable:
Run the following command from this directory:

```go build -o <executable-name>```
and then add your executable to your PATH

## To Run:
```<executable-name> 'Explain Big O notation.'```