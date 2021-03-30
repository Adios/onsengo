package cmd

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/adios/onsengo/onsen"
)

func init() {
	lsCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "list all episodes as well")
	rootCmd.AddCommand(lsCmd)
}

var (
	recursive bool
	lsCmd     = &cobra.Command{
		Use:   "ls",
		Short: "List radio shows",
		Run:   lsRun,
	}
)

func lsRun(cmd *cobra.Command, args []string) {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c := &http.Client{Transport: t}

	req, err := http.NewRequest("GET", backend, nil)
	if err != nil {
		panic(err)
	}

	if session != "" {
		req.Header.Add("Cookie", "_session_id="+session)
	}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	o, err := onsen.Create(string(body))
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		rs := o.RadioShows()

		sort.SliceStable(rs, func(i, j int) bool {
			a, _ := rs[i].JstUpdatedAt()
			b, _ := rs[j].JstUpdatedAt()
			return a.Before(b)
		})

		if recursive {
			for _, r := range rs {
				t, _ := r.JstUpdatedAt()

				fmt.Printf(
					"d---- %3d %s %-20s %s\n",
					len(r.Episodes()),
					t.Format("Jan _2 2006"),
					r.Name(),
					r.Title(),
				)

				es := r.Episodes()
				for _, e := range es {
					var latest, playable, bonus, premium string

					et, _ := e.JstUpdatedAt()

					if e.IsLatest() {
						latest = "*"
					} else {
						latest = "-"
					}
					em, _ := e.Manifest()

					if em == "" {
						playable = "-"
					} else {
						playable = "r"
					}
					if e.IsBonus() {
						bonus = "b"
					} else {
						bonus = "-"
					}
					if e.RequiresPremium() {
						premium = "$"
					} else {
						premium = "-"
					}

					fmt.Printf(
						"-%s%s%s%s   1 %s %-20s %s\n", playable, bonus, latest, premium,
						et.Format("Jan _2 2006"),
						r.Name()+"/"+strconv.FormatUint(uint64(e.EpisodeId()), 10),
						e.Title(),
					)
				}
			}
		} else {
			longestEpisodes := 0
			longestName := 0

			for _, r := range rs {
				s := len(strconv.Itoa(len(r.Episodes())))
				if s > longestEpisodes {
					longestEpisodes = s
				}
				s = len(r.Name())
				if s > longestName {
					longestName = s
				}
			}

			for _, r := range rs {
				t, _ := r.JstUpdatedAt()

				fmt.Printf(
					"d---- %*d %s %-*s %s\n",
					longestEpisodes,
					len(r.Episodes()),
					t.Format("Jan _2 2006"),
					longestName,
					r.Name(),
					r.Title(),
				)
			}
		}
	} else {
		for _, arg := range args {
			v, ok := o.RadioShow(onsen.RadioShowName(arg))
			if ok == false {
				panic(arg)
			}

			es := v.Episodes()

			longestName := 0

			for _, e := range es {
				s := len(strconv.FormatUint(uint64(e.EpisodeId()), 10))
				if s > longestName {
					longestName = s
				}
			}

			longestName += len(arg) + 1

			for _, e := range es {
				var latest, playable, bonus, premium string

				if e.IsLatest() {
					latest = "*"
				} else {
					latest = "-"
				}

				em, _ := e.Manifest()

				if em == "" {
					playable = "-"
				} else {
					playable = "r"
				}
				if e.IsBonus() {
					bonus = "b"
				} else {
					bonus = "-"
				}
				if e.RequiresPremium() {
					premium = "$"
				} else {
					premium = "-"
				}

				et, _ := e.JstUpdatedAt()

				fmt.Printf(
					"-%s%s%s%s 1 %s %-*s %s\n", playable, bonus, latest, premium,
					et.Format("Jan _2 2006"),
					longestName,
					arg+"/"+strconv.FormatUint(uint64(e.EpisodeId()), 10),
					e.Title(),
				)

			}
		}
	}
}
