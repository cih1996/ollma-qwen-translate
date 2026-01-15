package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"
)

type Translator struct {
	ollamaURL string
	modelName string
	client    *http.Client
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewTranslator(ollamaURL, modelName string) *Translator {
	return &Translator{
		ollamaURL: ollamaURL,
		modelName: modelName,
		client:    &http.Client{},
	}
}

// DetectLanguage 检测文本语言（中文或英文）
func (t *Translator) DetectLanguage(text string) string {
	chineseCount := 0
	totalChars := 0

	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			chineseCount++
		}
		if unicode.IsLetter(r) || unicode.Is(unicode.Han, r) {
			totalChars++
		}
	}

	// 如果中文字符占比超过30%，认为是中文
	if totalChars > 0 && float64(chineseCount)/float64(totalChars) > 0.3 {
		return "chinese"
	}
	return "english"
}

// Translate 翻译文本
func (t *Translator) Translate(text string) (string, error) {
	// 检测语言
	lang := t.DetectLanguage(text)
	
	// 构建提示词
	var prompt string
	if lang == "chinese" {
		prompt = fmt.Sprintf("请将以下中文翻译成英文，只输出翻译结果，不要添加任何解释：\n%s", text)
	} else {
		prompt = fmt.Sprintf("Please translate the following English text to Chinese, only output the translation result without any explanation:\n%s", text)
	}

	// 调用 Ollama API
	result, err := t.callOllama(prompt)
	if err != nil {
		return "", fmt.Errorf("调用 Ollama API 失败: %w", err)
	}

	return strings.TrimSpace(result), nil
}

// callOllama 调用 Ollama API
func (t *Translator) callOllama(prompt string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", t.ollamaURL)

	reqBody := OllamaRequest{
		Model:  t.modelName,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API 返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Response, nil
}
