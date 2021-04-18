package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/adios/onsengo/onsen"
	"github.com/spf13/cobra"
)

var lsm = struct {
	after JstHyphenDate

	cmd *cobra.Command
}{
	cmd: &cobra.Command{
		Use:   "lsm [radio_name...] [radio_name/episode_id...]",
		Short: "List episode's manifest",
		Long: `
List all the episode manifests of each radio show. Pass a radio name to list 
only the episodes under the radio show. Pass a name-id format to show only the
specified manifest.

  onsengo lsm fujita             # list all manifests for a radio show
  onsengo lsm fujita/3919        # show specified episode
  onsengo lsm --after 2020-12-27 # list those updated on or after 2020/12/27 in JST

Note that inaccessible episodes are not shown. 
`,
	},
}

func init() {
	root.cmd.AddCommand(lsm.cmd)

	lsm.cmd.RunE = runLsm
	lsm.cmd.Flags().Var(&lsm.after, "after", "show only those that are on or after this date (in JST)")
}

func runLsm(cmd *cobra.Command, args []string) error {
	o, err := root.onsen()
	if err != nil {
		return err
	}

	// filtering on these cases: empty manifest, on-air before a given date
	f := NewFilter()
	if lsm.after != (JstHyphenDate{}) {
		f.With(FilterUpdatedAfter(time.Time(lsm.after)))
	}

	switch n := len(args); {
	case n == 0:
		o.EachRadio(func(r onsen.Radio) {
			for _, e := range r.Episodes() {
				f.Push(e)
			}
		})
	case n > 0:
		for _, arg := range unique(args) {
			m := regexp.MustCompile("[^/]+/([0-9]+)").FindStringSubmatch(arg)

			switch len(m) > 1 {
			case true:
				id, _ := strconv.Atoi(m[1])
				processDesignatedEpisode(o, f, id, arg)
			case false:
				r, ok := o.Radio(arg)
				if !ok {
					fmt.Fprintf(root.errw(), "%s: not found\n", arg)
					continue
				}
				for _, e := range r.Episodes() {
					f.Push(e)
				}
			}
		}
	}

	out := root.outw()
	for _, m := range f.Out() {
		fmt.Fprintf(out, "%s\n", m)
	}

	return nil
}

func processDesignatedEpisode(o *onsen.Onsen, f *Filter, id int, name string) {
	e, ok := o.Episode(id)
	if !ok {
		fmt.Fprintf(root.errw(), "%s: not found\n", name)
		return
	}
	_, ok = e.Manifest()
	if !ok {
		// Empty manifest on a designated episode should trigger a warning.
		fmt.Fprintf(root.errw(), "%s: empty manifest, may be inaccessible\n", name)
		return
	}
	f.Push(e)
}

type (
	// Filter stores a chain of if-else procedure and run these tests on each Push() to filter
	// input episodes. It stores all the manifests of the episodes that passed the filtering.
	Filter struct {
		q     []string
		chain []FilterFn
	}

	FilterFn  func(Episoder) bool
	FilterOpt func(*Filter)

	Episoder interface {
		Id() int
		JstUpdatedAt() (time.Time, bool)
		Manifest() (string, bool)
	}
)

// Takes an episode to run through the filter chain.
func (f *Filter) Push(e Episoder) {
	for _, pass := range f.chain {
		if !pass(e) {
			return
		}
	}

	m, ok := e.Manifest()
	if !ok {
		return
	}

	f.q = append(f.q, m)
}

func (f *Filter) Out() []string {
	return f.q
}

func (f *Filter) With(opts ...FilterOpt) {
	for _, opt := range opts {
		opt(f)
	}
}

func NewFilter(opts ...FilterOpt) *Filter {
	f := &Filter{
		q:     []string{},
		chain: []FilterFn{},
	}
	f.With(opts...)

	return f
}

// Filters out an episode if it was updated before the given date.
func FilterUpdatedAfter(dt time.Time) FilterOpt {
	fn := func(e Episoder) bool {
		tm, ok := e.JstUpdatedAt()
		if !ok {
			log.Printf("%d: doesn't have update time, filtered\n", e.Id())
			return false
		}
		return !tm.Before(dt)
	}

	return func(f *Filter) {
		f.chain = append(f.chain, fn)
	}
}

// A custom date format to fulfill pflag.Value interface.
type JstHyphenDate time.Time

func (h *JstHyphenDate) Set(dt string) error {
	loc := time.FixedZone("UTC+9", 9*60*60)

	tm, err := time.ParseInLocation("2006-01-02", dt, loc)
	if err != nil {
		return err
	}

	*h = JstHyphenDate(tm)

	return nil
}

func (h *JstHyphenDate) Type() string {
	return "YYYY-MM-DD"
}

func (h *JstHyphenDate) String() string {
	return ""
}
