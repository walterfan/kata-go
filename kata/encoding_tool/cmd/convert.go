package cmd

import (
	"encoding/base64"
	"net/url"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert between different encodings",
	Long:  `Convert between various encodings like base64, URL encoding, etc.`,
}

var base64EncodeCmd = &cobra.Command{
	Use:   "base64encode [text]",
	Short: "Encode text to base64",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encoded := base64.StdEncoding.EncodeToString([]byte(args[0]))
		cmd.Println(encoded)
	},
}

var base64DecodeCmd = &cobra.Command{
	Use:   "base64decode [text]",
	Short: "Decode base64 text",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		decoded, err := base64.StdEncoding.DecodeString(args[0])
		if err != nil {
			cmd.Println("Error decoding:", err)
			return
		}
		cmd.Println(string(decoded))
	},
}

var urlEncodeCmd = &cobra.Command{
	Use:   "urlencode [text]",
	Short: "URL encode text",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encoded := url.QueryEscape(args[0])
		cmd.Println(encoded)
	},
}

var urlDecodeCmd = &cobra.Command{
	Use:   "urldecode [text]",
	Short: "URL decode text",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		decoded, err := url.QueryUnescape(args[0])
		if err != nil {
			cmd.Println("Error decoding:", err)
			return
		}
		cmd.Println(decoded)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.AddCommand(base64EncodeCmd)
	convertCmd.AddCommand(base64DecodeCmd)
	convertCmd.AddCommand(urlEncodeCmd)
	convertCmd.AddCommand(urlDecodeCmd)
}
