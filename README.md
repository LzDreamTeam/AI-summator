# AI Summator

This is a CLI application written in Go that uses a local LLM (llama3.1 via Ollama) to sum two numbers. It demonstrates how to integrate Go with local LLMs using LangChain for Go (`langchaingo`).

## Prerequisites

1.  **Go**: Ensure you have Go installed (version 1.22 or later recommended).
2.  **Ollama**: You need to have [Ollama](https://ollama.com/) installed and running.
3.  **Llama 3.1 Model**: Pull the required model:
    ```bash
    ollama pull llama3.1
    ```

## Installation

1.  Clone the repository (if applicable) or navigate to the project directory.
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Build the application:
    ```bash
    go build -o ai-summator main.go
    ```

## Usage

Run the built binary with two numeric arguments:

```bash
./ai-summator 5 3
```

Example output:
```
Result: 8.000000
```

Floating point numbers are supported:

```bash
./ai-summator 1.5 2.7
```

## Testing

The project includes both unit tests and integration tests.

To run all tests:

```bash
go test -v ./...
```

**Note**: The integration tests require Ollama to be running and the `llama3.1` model to be available. If Ollama is not reachable, the integration test will fail.

## Project Structure

-   `main.go`: Entry point for the CLI.
-   `summator/`: Contains the core logic and tests.
    -   `summator.go`: Implementation of the summator using `langchaingo`.
    -   `summator_test.go`: Unit tests with mocked LLM.
    -   `integration_test.go`: Integration tests against a real Ollama instance.
