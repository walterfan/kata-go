package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/walterfan/blog-gen/internal"
	"gopkg.in/yaml.v3"
)

var idea string
var location string
var model string
var titleArg string
var language string

// AppConfig defines YAML structure for prompts and location
type AppConfig struct {
	Prompts struct {
		System string `yaml:"system"`
		User   string `yaml:"user"`
	} `yaml:"prompts"`
	Location string `yaml:"location"`
}

func loadAppConfig(path string) (*AppConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg AppConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

var rootCmd = &cobra.Command{
	Use:   "bloggen",
	Short: "Generate a technical blog with weather information using OpenAI API",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			logrus.Warn("No .env file found")
		}

		// Load config.yaml if present
		var fileCfg *AppConfig
		if cfg, err := loadAppConfig("config/config.yaml"); err == nil {
			fileCfg = cfg
		} else {
			logrus.Debugf("No config loaded or parse error: %v", err)
		}

		if idea == "" {
			idea = "用 WebRTC 和 Pion 打造一款网络录音机"
		}

		// Decide location: CLI flag overrides config
		if fileCfg != nil && fileCfg.Location != "" && (location == "" || location == "Beijing") {
			location = fileCfg.Location
		}
		if location == "" {
			location = "Hefei"
		}

		today := time.Now().Format("2006-01-02")

		// Compute title (CLI overrides default)
		computedTitle := titleArg
		if strings.TrimSpace(computedTitle) == "" {
			computedTitle = fmt.Sprintf("my blog at %s", today)
		}

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
			"title":    computedTitle,
			"idea":     idea,
			"location": location,
		})
		if err != nil {
			logrus.Fatalf("Render error: %v", err)
		}

		// Build prompts from config with fallback
		systemPrompt := "你是一个技术科普作家和资深的内容创作者, 行文风趣幽默, 发人深省"
		if fileCfg != nil && strings.TrimSpace(fileCfg.Prompts.System) != "" {
			systemPrompt = fileCfg.Prompts.System
		}

		userPrompt := fmt.Sprintf(`
我有一个技术博客, 用来分享自己在技术上的想法和心得, 请根据如下模板为我生成今天的博客内容, 替换掉模板中的 "..." 字符串
--------------
%s
请使用 Markdown 格式分别输出英文和中文的博客内容，并在天气部分填入今天的天气信息`, rendered)
		if fileCfg != nil && strings.TrimSpace(fileCfg.Prompts.User) != "" {
			userPrompt = strings.ReplaceAll(fileCfg.Prompts.User, "{{ template }}", rendered)
		}

		// Load LLM configuration
		cfg := internal.LoadLlmConfigFromEnv()
		if model != "" {
			cfg.Model = model
		}
		service := internal.NewLlmService(cfg)

		// Generate blog content
		content, err := service.FullToolAwareChatFlow(cmd.Context(), systemPrompt, userPrompt, location)
		if err != nil {
			logrus.Fatalf("LLM call failed: %v", err)
		}

		// Split into English and Chinese parts
		enPart, zhPart := splitEnglishChinese(content)
		if strings.TrimSpace(enPart) == "" && strings.TrimSpace(zhPart) == "" {
			logrus.Warn("Generated content is empty; nothing to write")
			return
		}

		//create output directory if not exists
		if _, err := os.Stat("output"); os.IsNotExist(err) {
			os.Mkdir("output", 0755)
		}
		// Write files, honoring language selection (en, zh, both or empty)
		filenameEn := fmt.Sprintf("output/blog-%s-en.md", today)
		filenameZh := fmt.Sprintf("output/blog-%s-cn.md", today)

		wantEn := language == "" || strings.EqualFold(language, "both") || strings.EqualFold(language, "en")
		wantZh := language == "" || strings.EqualFold(language, "both") || strings.EqualFold(language, "zh")

		if wantEn && strings.TrimSpace(enPart) != "" {
			if err := os.WriteFile(filenameEn, []byte(enPart), 0644); err != nil {
				logrus.Fatalf("Save English blog failed: %v", err)
			}
			logrus.Infof("English blog saved to %s", filenameEn)
		}

		if wantZh && strings.TrimSpace(zhPart) != "" {
			if err := os.WriteFile(filenameZh, []byte(zhPart), 0644); err != nil {
				logrus.Fatalf("Save Chinese blog failed: %v", err)
			}
			logrus.Infof("Chinese blog saved to %s", filenameZh)
		}
	},
}

func splitEnglishChinese(content string) (string, string) {

	// find a top-level header line with any non-ASCII (likely Chinese)
	reZhHeader := regexp.MustCompile(`(?m)^#\s+.*[^\x00-\x7F].*$`)
	if m := reZhHeader.FindStringIndex(content); m != nil {
		en := strings.TrimSpace(content[:m[0]])
		zh := strings.TrimSpace(content[m[0]:])
		return en, zh
	}
	// If not found, try splitting at the second top-level header
	reAllHeaders := regexp.MustCompile(`(?m)^#\s+.*$`)
	matches := reAllHeaders.FindAllStringIndex(content, -1)
	if len(matches) >= 2 {
		en := strings.TrimSpace(content[:matches[1][0]])
		zh := strings.TrimSpace(content[matches[1][0]:])
		return en, zh
	}
	// As last resort, return whole content as English
	return strings.TrimSpace(content), ""
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&idea, "idea", "i", "", "Today's technical inspiration/title")
	rootCmd.PersistentFlags().StringVarP(&titleArg, "title", "t", "", "The blog title (overrides template title)")
	rootCmd.PersistentFlags().StringVarP(&language, "language", "l", "", "Blog language: en, zh, or both")
	rootCmd.PersistentFlags().StringVarP(&location, "city", "c", "Beijing", "City for weather information")
	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", "", "OpenAI model name (e.g., gpt-4o)")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
