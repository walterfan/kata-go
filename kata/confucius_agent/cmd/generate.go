package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate random values",
	Long:  `Generate random values like UUIDs, random strings, etc.`,
}

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate a UUID",
	Run: func(cmd *cobra.Command, args []string) {
		id := uuid.New()
		cmd.Println(id.String())
	},
}

var randomStringCmd = &cobra.Command{
	Use:   "random [length]",
	Short: "Generate a random string",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		length := 16 // default length
		if len(args) > 0 {
			_, err := fmt.Sscanf(args[0], "%d", &length)
			if err != nil {
				cmd.Println("Invalid length:", args[0])
				return
			}
		}

		withNumbers, _ := cmd.Flags().GetBool("numbers")
		withSymbols, _ := cmd.Flags().GetBool("symbols")

		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		if withNumbers {
			charset += "0123456789"
		}
		if withSymbols {
			charset += "!@#$%^&*()_+-=[]{}|;:,.<>?"
		}

		result := make([]byte, length)
		for i := 0; i < length; i++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				cmd.Println("Error generating random number:", err)
				return
			}
			result[i] = charset[num.Int64()]
		}

		cmd.Println(string(result))
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.AddCommand(uuidCmd)
	generateCmd.AddCommand(randomStringCmd)

	randomStringCmd.Flags().BoolP("numbers", "n", false, "Include numbers")
	randomStringCmd.Flags().BoolP("symbols", "s", false, "Include symbols")
}
