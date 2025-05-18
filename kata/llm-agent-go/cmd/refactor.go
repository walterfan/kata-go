package cmd

import (
	"fmt"
	"github.com/walterfan/llm-agent-go/internal/llm"
	"github.com/walterfan/llm-agent-go/internal/prompt"
	"os"

	"github.com/spf13/cobra"
)

var refactorCmd = &cobra.Command{
	Use:   "refactor [file]",
	Short: "Refactor Go source code using LLM",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		code, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		p := prompt.BuildRefactorPrompt(string(code))
		resp, err := llm.AskLLM(p)
		if err != nil {
			fmt.Println("LLM error:", err)
			return
		}
		fmt.Println("🛠 Refactored Code:\n", resp)
	},
}

func init() {
	rootCmd.AddCommand(refactorCmd)
}
