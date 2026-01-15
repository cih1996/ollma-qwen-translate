# Project Overview: Ollama Qwen Translate

## Introduction
`ollma-qwen-translate` (CLI alias: `fy`) is a Go-based command-line tool designed to perform intelligent text translation between Chinese and English. It leverages a local [Ollama](https://ollama.ai/) instance to run Large Language Models (LLMs) for high-quality translation.

## Project Structure
(Excluding `VibeVoice` directory)

```
.
├── cmd/
│   └── root.go           # CLI command definition, flags, and execution entry point
├── internal/
│   └── translator/
│       └── translator.go # Core translation logic, language detection, and Ollama API integration
├── main.go               # Application entry point
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── README.md             # User documentation
```

## Code Analysis

### Entry Point (`main.go`)
The application starts in `main.go`, which simply delegates execution to `cmd.Execute()`. This pattern keeps the `main` package clean and allows the `cmd` package to handle CLI logic.

### CLI Logic (`cmd/root.go`)
- **Library**: Uses `github.com/spf13/cobra` for CLI command management.
- **Command**: `fy`
- **Arguments**: Accepts one argument (text to translate) OR zero arguments if `-c` is used.
- **Flags**:
  - `-u, --url`: Ollama server URL (Default: `http://localhost:11434`)
  - `-m, --model`: Model name (Default: `huihui_ai/qwen3-abliterated:0.6b`)
  - `-c, --clipboard`: Read input text from system clipboard.
- **Execution Flow**:
  1. Parses flags and arguments.
  2. Initializes a `Translator` instance.
  3. Calls `trans.Translate(inputText)`.
  4. Prints the result or exits on error.

### Core Logic (`internal/translator/translator.go`)

#### 1. Language Detection (`DetectLanguage`)
The tool uses a heuristic approach to determine if the input is Chinese or English:
- Iterates through the input string.
- Counts Chinese characters (`unicode.Han`) and total letters/Chinese characters.
- **Rule**: If the ratio of Chinese characters to total significant characters exceeds **30%**, it is classified as **Chinese**. Otherwise, it is classified as **English**.

#### 2. Prompt Engineering (`Translate`)
Based on the detected language, a specific prompt is constructed to instruct the LLM:
- **Chinese Input**: "请将以下中文翻译成英文，只输出翻译结果，不要添加任何解释：\n{text}"
- **English Input**: "Please translate the following English text to Chinese, only output the translation result without any explanation:\n{text}"

#### 3. Ollama Integration (`callOllama`)
- **Endpoint**: Sends a POST request to `{ollamaURL}/api/generate`.
- **Payload**:
  ```json
  {
    "model": "huihui_ai/qwen3-abliterated:0.6b",
    "prompt": "...",
    "stream": false
  }
  ```
- **Response**: Parses the JSON response to extract the generated text.

## Dependencies
- **Cobra**: For creating powerful modern CLI applications.
- **Standard Library**: heavily uses `net/http`, `encoding/json`, `unicode`, etc.

## Configuration
The tool defaults to using `huihui_ai/qwen3-abliterated:0.6b`, which implies this specific model should be pulled in Ollama before use:
```bash
ollama pull huihui_ai/qwen3-abliterated:0.6b
```
