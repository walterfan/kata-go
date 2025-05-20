package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/walterfan/llm-agent-go/internal/llm"
	"github.com/walterfan/llm-agent-go/internal/log"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

type PromptConfig struct {
	Name         string `mapstructure:"name" json:"name"`
	Description  string `mapstructure:"description" json:"description"`
	SystemPrompt string `mapstructure:"system_prompt" json:"system_prompt"`
	UserPrompt   string `mapstructure:"user_prompt" json:"user_prompt"`
	Tags         string `mapstructure:"tags" json:"tags"`
}

func getPromptConfig(commandName string) (*PromptConfig, error) {
	promptList := viper.Get("prompts").([]interface{})
	for _, p := range promptList {
		itemMap := p.(map[string]interface{})
		if itemMap["name"].(string) == commandName {
			return &PromptConfig{
				Name:         itemMap["name"].(string),
				Description:  itemMap["description"].(string),
				SystemPrompt: itemMap["system_prompt"].(string),
				UserPrompt:   itemMap["user_prompt"].(string),
				Tags:         itemMap["tags"].(string),
			}, nil
		}
	}
	return nil, fmt.Errorf("prompt not found for command: %s", commandName)
}

func processCommand(cmd *cobra.Command, args []string) {
	path := args[0]
	commandName, _ := cmd.Flags().GetString("command")
	streamMode, _ := cmd.Flags().GetBool("stream")
	outputLanguage, _ := cmd.Flags().GetString("language")

	code, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	selectedPrompt, err := getPromptConfig(commandName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Build final prompt by replacing {{code}}
	promptText := strings.Replace(selectedPrompt.UserPrompt, "{{code}}", string(code), -1)

	if outputLanguage != "" {
		promptText = strings.Replace(promptText, "{{output_language}}", outputLanguage, -1)
	}

	if streamMode {
		err = llm.AskLLMWithStream(selectedPrompt.SystemPrompt, promptText, func(chunk string) {
			fmt.Print(chunk)
		})
	} else {
		resp, err := llm.AskLLM(selectedPrompt.SystemPrompt, promptText)
		if err == nil {
			fmt.Printf("📝 Result (%s):\n%s\n", commandName, resp)
		}
	}

	if err != nil {
		fmt.Println("LLM error:", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "coder-helper  <file>",
	Short: "An LLM-powered Go code assistant for explain, review, refactor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		processCommand(cmd, args)
	},
}

func init() {
	var err error
	logger, err = log.InitLogger()
	if err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}

	rootCmd.Flags().StringP("command", "c", "review", "Specify the command (e.g. explain, review, refactor)")
	rootCmd.Flags().StringP("language", "l", "golang", "Specify the output language (e.g. golang, java, python, etc.)")
	rootCmd.Flags().StringP("tongue", "t", "chinese", "Specify the output language (e.g. chinese, english, japanese, etc.)")
	rootCmd.Flags().BoolP("stream", "s", false, "Enable streaming mode for LLM response")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
