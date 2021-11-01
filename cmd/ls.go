package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/adios/onsengo/onsen"
	pp "github.com/adios/pprint"
)

var ls = struct {
	recursive bool

	cmd *cobra.Command

	// map for translating boolean to ls letters
	lut letters
}{
	cmd: &cobra.Command{
		Use:   "ls [radio_name...]",
		Short: "List radio shows",
		Long: `
List and browse radio shows, outputs in ascending order on uploaded/published 
date. Provide radio names to list only those shows including their episodes.

Use -r to list all radio shows and their episodes.
`,
	},
}

func init() {
	root.cmd.AddCommand(ls.cmd)

	ls.cmd.RunE = runLs
	ls.cmd.Flags().BoolVarP(&ls.recursive, "recursive", "r", false, "include all episodes")
}

func runLs(cmd *cobra.Command, args []string) error {
	o, err := root.onsen()
	if err != nil {
		return err
	}

	setupLs()

	out := typeset()

	switch n := len(args); {
	case n == 0:
		if ls.recursive {
			o.EachRadio(func(r onsen.Radio) { addRadioEpisodes(out, r) })
		} else {
			o.EachRadio(func(r onsen.Radio) { addRadio(out, r) })
		}
	case n > 0:
		for _, arg := range unique(args) {
			r, ok := o.Radio(arg)
			if !ok {
				fmt.Fprintf(root.errw(), "%s: not found\n", arg)
				continue
			}
			addRadioEpisodes(out, r)
		}
	}

	// 0, 1, 2: letters, count, mtime ...
	if err := out.Sort(2, pp.WithCmpMatchers(mtimeCmp)); err != nil {
		return err
	}

	pp.Print(out, pp.WithWriter(root.outw()))

	return nil
}

// Returns the pushed node to create folder-like context to further push episodes to it.
//
// output
//   - radio name
//     - episode
//     - episode
//     - ...
//   - radio name
//   - ...
//
// Sort on output (root level) affects only on "radio name" level.
func addRadio(out *pp.Node, r onsen.Radio) (pushed *pp.Node) {
	tm, _ := r.JstUpdatedAt()

	pushed, _ = out.Push(
		toRadioLetters(r),
		len(r.Episodes()),
		mtime(tm),
		r.Name(),
		r.Title(),
	)

	return pushed
}

func addRadioEpisodes(out *pp.Node, r onsen.Radio) {
	var (
		// Push radio first
		dir     = addRadio(out, r)
		dirName = r.Name()
	)

	// And then push every episodes under that radio
	for _, e := range r.Episodes() {
		tm, _ := e.JstUpdatedAt()

		// Append guests to radio episode title
		last := e.Title()
		if len(e.Guests()) != 0 {
			last += " #"

			for _, p := range e.Guests() {
				last += " " + p.Name()
			}
		}

		dir.Push(
			toEpisodeLetters(e),
			1,
			mtime(tm),
			dirName+"/"+strconv.FormatInt(int64(e.Id()), 10),
			last,
		)
	}
}

func toRadioLetters(r onsen.Radio) string {
	return "d--" + ls.lut["just updated"][r.HasBeenUpdated()] + "--"
}

func toEpisodeLetters(e onsen.Episode) string {
	var (
		m    = ls.lut
		u, _ = e.Manifest()
	)

	return "-" +
		m["accessible"][u != ""] +
		m["include video"][e.HasVideoStream()] +
		m["just updated"][e.IsLatest()] +
		m["extra content"][e.IsBonus()] +
		m["paid content"][e.RequiresPremium()]
}

func typeset() *pp.Node {
	return pp.NewNode(
		pp.WithColumns(
			pp.NewColumn(),                       // radio / episode letters
			pp.NewColumn(),                       // episodes
			pp.NewColumn(),                       // mtime
			pp.NewColumn(pp.WithLeftAlignment()), // name
			pp.NewColumn(pp.WithWidth(0)),        // title / title + guests
		),
	)
}

// Custom-defined time for String() to pprint's formatter.
type mtime time.Time

func (m mtime) String() string {
	return time.Time(m).Format("Jan _2 2006")
}

// To be able to sort on our custom type.
func mtimeCmp(a interface{}) pp.CmpFn {
	return func(a, b interface{}) bool {
		return time.Time(a.(mtime)).Before(time.Time(b.(mtime)))
	}
}

type letters map[string]map[bool]string

func setupLs() {
	ls.lut = letters{
		"accessible": {
			true:  "r",
			false: "-",
		},
		"include video": {
			true:  "v",
			false: "-",
		},
		"just updated": {
			true:  "*",
			false: "-",
		},
		"extra content": {
			true:  "+",
			false: "-",
		},
		"paid content": {
			true:  "$",
			false: "-",
		},
	}

}

// Stable unique.
func unique(s []string) []string {
	var out []string

	switch n := len(s); {
	case n < 2:
		out = s
	case n == 2:
		if s[0] == s[1] {
			out = s[0:1]
		} else {
			out = s
		}
	default:
		out = make([]string, 0, n)

		bucket := make(map[string]bool)
		for _, key := range s {
			if bucket[key] == false {
				bucket[key] = true
				out = append(out, key)
			}
		}
	}
	return out
}
