package cmd

import (
	"github.com/spf13/cobra"
)

// Use for generating resources
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Use for generating utilities",
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
