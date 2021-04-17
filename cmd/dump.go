package cmd

import (
	"fmt"

	"github.com/adios/onsengo/onsen"
	"github.com/spf13/cobra"
)

var dump = struct {
	cmd *cobra.Command
}{
	cmd: &cobra.Command{
		Use:   "dump",
		Short: "Dump raw data from Onsen website into a json",
		Long: `
All the Onsen's website data is stored in a minified Javascript object. This
command requests to the website, parses the content, and unminifys the object,
providing a way to access the raw data in json format. You can then pass the
output to "jq" to process the json manually.
`,
		RunE: runDump,
	},
}

func init() {
	root.cmd.AddCommand(dump.cmd)
}

func runDump(cmd *cobra.Command, args []string) error {
	html, err := root.html()
	if err != nil {
		return err
	}

	str, err := onsen.RawData(html)
	if err != nil {
		return err
	}

	fmt.Fprintf(root.outw(), "%s\n", str)

	return nil
}
