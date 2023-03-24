package main

import (
	"os"

	"migraterr/cmd"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "migraterr",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(cmd.RunBencode())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
