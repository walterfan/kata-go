package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/walterfan/blog-gen/internal"
)

var idea string
var location string
var model string

var rootCmd = &cobra.Command{
	Use:   "bloggen",
	Short: "Generate a technical blog with weather information using OpenAI API",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			logrus.Warn("No .env file found")
		}

		if idea == "" {
			idea = "用 WebRTC 和 Pion 打造一款网络录音机"
		}
		if location == "" {
			location = "Hefei"
		}

		today := time.Now().Format("2006-01-02")

		// Read template
		templateStr, err := os.ReadFile("templates/blog.md.teml")
		if err != nil {
			logrus.Fatalf("Read template error: %v", err)
		}

		tmpl, err := pongo2.FromString(string(templateStr))
		if err != nil {
			logrus.Fatalf("Parse template error: %v", err)
		}

		// Render template with basic data
		rendered, err := tmpl.Execute(pongo2.Context{
			"title": fmt.Sprintf("my blog at %s", today),
			"idea":  idea,
		})
		if err != nil {
			logrus.Fatalf("Render error: %v", err)
		}

		// Create system and user prompts similar to Python version
		systemPrompt := "你是一个技术科普作家和资深的内容创作者, 行文风趣幽默, 发人深省"
		userPrompt := fmt.Sprintf(`
我有一个技术博客, 用来分享自己在技术上的想法和心得, 请根据如下模板为我生成今天的博客内容, 替换掉模板中的 "..." 字符串
--------------
%s
请使用 Markdown 格式分别输出英文和中文的博客内容，并在天气部分填入今天的天气信息`, rendered)

		// Load LLM configuration
		cfg := internal.LoadLlmConfigFromEnv()
		if model != "" {
			cfg.Model = model
		}
		service := internal.NewLlmService(cfg)

		// Generate blog content (using simple version for debugging)
		content, err := service.SimpleFullToolAwareChatFlow(cmd.Context(), systemPrompt, userPrompt, location)
		if err != nil {
			logrus.Fatalf("LLM call failed: %v", err)
		}

		// Save blog to file
		filename := fmt.Sprintf("blog-%s.md", today)
		err = os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			logrus.Fatalf("Save blog failed: %v", err)
		}
		logrus.Infof("Blog saved to %s", filename)
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&idea, "idea", "i", "", "Today's technical inspiration/title")
	rootCmd.PersistentFlags().StringVarP(&location, "location", "l", "Beijing", "City for weather information")
	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", "", "OpenAI model name (e.g., gpt-4o)")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
