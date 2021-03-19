package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "onsengo",
	Short: "List and download radio programs from onsen.ag.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
