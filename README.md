# fy - Smart Translation Tool

A smart translation CLI tool based on Ollama, which automatically detects Chinese and English and translates between them.

## Features

- 🤖 Automatically detects input language (Chinese/English)
- 🔄 Automatically translates to the target language
- 📋 Supports reading from clipboard
- ⚡ Fast response based on local Ollama service
- 📦 Can be packaged as a standalone CLI tool

## Prerequisites

1. Go 1.21 or higher installed
2. Ollama service deployed and running (Default address: http://localhost:11434)
3. Model downloaded and configured (Default uses `huihui_ai/qwen3-abliterated:0.6b`, configurable via `-m` parameter)

## Installation

### Build from Source

```bash
# Clone or download the project
git clone <repository-url>
cd ollma-qwen-translate/Translate

# Install dependencies
go mod download

# Build
go build -o fy

# Move the executable to PATH (Optional)
sudo mv fy /usr/local/bin/
```

### Cross Compilation

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o fy-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o fy-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o fy.exe
```

## Usage

### Basic Usage

```bash
# Translate Chinese to English
fy "你好"

# Translate English to Chinese
fy "hello"

# Translate long text
fy "你好，世界！这是一个测试。"

# Read from clipboard and translate
fy -c
```

### Advanced Options

```bash
# Specify Ollama server address
fy -u http://localhost:11434 "你好"

# Specify model
fy -m qwen2.5 "hello"

# Combined usage
fy -u http://localhost:11434 -m qwen2.5 "你好世界"
```

## Parameters

- `-u, --url`: Ollama server address (Default: http://localhost:11434)
- `-m, --model`: Model name (Default: `huihui_ai/qwen3-abliterated:0.6b`)
- `-c, --clipboard`: Read content from clipboard and translate

## Language Detection Logic

The tool uses a simple heuristic method to detect language:
- If the ratio of Chinese characters in the text exceeds 30%, it is identified as Chinese and translated to English.
- Otherwise, it is identified as English and translated to Chinese.

## Troubleshooting

### Ollama Connection Failure

Ensure Ollama service is running:
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags
```

### Model Not Found

Ensure the specified model is downloaded:
```bash
# List installed models
ollama list

# Download model (if needed)
ollama pull huihui_ai/qwen3-abliterated:0.6b
```

## Development

```bash
cd Translate

# Run tests
go test ./...

# Format code
go fmt ./...

# Run
go run main.go "你好"
```

## License

MIT
