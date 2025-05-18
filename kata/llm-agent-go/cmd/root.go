package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "llm-agent",
	Short: "An LLM-powered Go code reviewer and refactorer",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
