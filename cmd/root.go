package cmd

import (
	"github.com/spf13/cobra"
)

var (
	backend string
	session string

	rootCmd = &cobra.Command{
		Use:   "onsengo",
		Short: "List and download onsen.ag radio shows",
		Long: `
onsengo♨ is a program which allows browsing radio shows on https://onsen.ag.
`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	fs := rootCmd.PersistentFlags()

	fs.StringVar(&backend, "backend", "https://onsen.ag/", "set backend url, scheme can be file://")
	fs.StringVarP(&session, "session", "s", "", "set session")
}
