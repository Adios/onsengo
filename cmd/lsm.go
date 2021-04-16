package cmd

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/adios/onsengo/onsen"
	"github.com/spf13/cobra"
)

var lsm = struct {
	cmd *cobra.Command
}{
	cmd: &cobra.Command{
		Use:   "lsm [radio_name...] [radio_name/episode_id...]",
		Short: "List episode's manifest",
		Long: `
List all the episode manifests of each radio show. Pass a radio name to list 
only the episodes under the radio show. Pass a name-id format to show only the
specified manifest.

  onsengo lsm fujita      # list all manifests for a radio show
  onsengo lsm fujita/3919 # show specified episode

Note that inaccessible episodes are not shown.
`,
		RunE: runLsm,
	},
}

func init() {
	root.cmd.AddCommand(lsm.cmd)
}

func runLsm(cmd *cobra.Command, args []string) error {
	o, err := root.onsen()
	if err != nil {
		return err
	}

	out := root.outw()

	switch n := len(args); {
	case n == 0:
		o.EachRadio(func(r onsen.Radio) {
			printEpisodesOf(r, out)
		})
	case n > 0:
		for _, arg := range unique(args) {
			m := regexp.MustCompile("[^/]+/([0-9]+)").FindStringSubmatch(arg)

			switch len(m) > 1 {
			case true:
				id, _ := strconv.Atoi(m[1])
				printEpisode(o, out, id, arg)
			case false:
				r, ok := o.Radio(arg)
				if !ok {
					fmt.Fprintf(root.errw(), "%s: not found\n", arg)
					continue
				}
				printEpisodesOf(r, out)
			}
		}
	}

	return nil
}

func printEpisode(o *onsen.Onsen, out io.Writer, id int, name string) {
	e, ok := o.Episode(id)
	if !ok {
		fmt.Fprintf(root.errw(), "%s: not found\n", name)
		return
	}
	m, ok := e.Manifest()
	if !ok {
		fmt.Fprintf(root.errw(), "%s: empty manifest, may be inaccessible\n", name)
		return
	}
	fmt.Fprintf(out, "%s\n", m)
}

func printEpisodesOf(r onsen.Radio, out io.Writer) {
	for _, e := range r.Episodes() {
		m, ok := e.Manifest()
		if !ok {
			continue
		}
		fmt.Fprintf(out, "%s\n", m)
	}
}
