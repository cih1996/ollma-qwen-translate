package cmd

import (
	"fmt"
	"os"

	"ollma-qwen-translate/internal/translator"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var (
	ollamaURL    string
	modelName    string
	inputText    string
	useClipboard bool
)

var rootCmd = &cobra.Command{
	Use:   "fy",
	Short: "一个智能翻译工具，自动识别中英文并翻译",
	Long: `fy 是一个基于 Ollama 的智能翻译工具。
它会自动检测输入文本的语言（中文或英文），然后翻译成对应的语言。

示例:
  fy "你好"
  fy "hello"
  fy "你好世界"
  fy -c  (自动读取剪贴板内容)`,
	Args: func(cmd *cobra.Command, args []string) error {
		if useClipboard {
			if len(args) > 0 {
				return fmt.Errorf("accepts 0 args when using clipboard, received %d", len(args))
			}
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if useClipboard {
			inputText, err = clipboard.ReadAll()
			if err != nil {
				fmt.Fprintf(os.Stderr, "读取剪贴板失败: %v\n", err)
				os.Exit(1)
			}
			if inputText == "" {
				fmt.Fprintf(os.Stderr, "剪贴板为空\n")
				os.Exit(1)
			}
			fmt.Printf("从剪贴板读取内容: %s\n", inputText)
		} else {
			inputText = args[0]
		}

		// 创建翻译器
		trans := translator.NewTranslator(ollamaURL, modelName)

		// 执行翻译
		result, err := trans.Translate(inputText)
		if err != nil {
			fmt.Fprintf(os.Stderr, "翻译失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(result)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&ollamaURL, "url", "u", "http://localhost:11434", "Ollama 服务器地址")
	rootCmd.Flags().StringVarP(&modelName, "model", "m", "huihui_ai/qwen3-abliterated:0.6b", "使用的模型名称")
	rootCmd.Flags().BoolVarP(&useClipboard, "clipboard", "c", false, "从剪贴板读取内容")
}

func Execute() error {
	return rootCmd.Execute()
}
