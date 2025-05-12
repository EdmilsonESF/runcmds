package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "runcmds",
	Short: "A brief description of your application",
	Long:  `A longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		print("hello")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
