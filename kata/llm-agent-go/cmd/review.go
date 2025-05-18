// cmd/review.go
package cmd

import (
	"fmt"
	"github.com/walterfan/llm-agent-go/internal/llm"
	"github.com/walterfan/llm-agent-go/internal/prompt"
	"os"

	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review [file]",
	Short: "Review Go source code",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		streamMode, _ := cmd.Flags().GetBool("stream")

		code, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		promptText := prompt.BuildReviewPrompt(path, string(code))

		if streamMode {
			// Use streaming mode
			err = llm.AskLLMWithStream(promptText, func(chunk string) {
				fmt.Print(chunk) // Print each chunk as it comes
			})
		} else {
			// Use normal mode
			resp, err := llm.AskLLM(promptText)
			if err == nil {
				fmt.Println("🔍 Review Report:\n", resp)
			}
		}

		if err != nil {
			fmt.Println("LLM error:", err)
		}
	},
}

func init() {
	reviewCmd.Flags().BoolP("stream", "s", false, "Enable streaming mode for LLM response")
	rootCmd.AddCommand(reviewCmd)
}