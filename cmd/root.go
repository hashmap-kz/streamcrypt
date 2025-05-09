package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "streamcrypt",
	Short:        "A CLI tool for compressing and encrypting files",
	Long:         `streamcrypt is a flexible tool for compression and encryption of file streams using a pipeline-based architecture.`,
	SilenceUsage: true,
}

// Call this from main
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
}
