package cmd

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/adios/onsengo/onsen"
	"github.com/adios/onsengo/onsen/adapter"
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

	onsen, err := onsen.Create(string(body))
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		rs := onsen.RadioShows()

		sort.SliceStable(rs, func(i, j int) bool {
			return rs[i].GuessedUpdatedAt().Before(rs[j].GuessedUpdatedAt())
		})

		if recursive {
			for _, r := range rs {
				fmt.Printf(
					"d---- %3d %s %-20s %s\n",
					len(r.Episodes()),
					r.GuessedUpdatedAt().Format("Jan _2 2006"),
					r.Name(),
					r.Title(),
				)

				es := r.Episodes()
				for _, e := range es {
					var latest, playable, bonus, premium string

					if e.IsLatest() {
						latest = "*"
					} else {
						latest = "-"
					}
					if e.Manifest() == "" {
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
						e.GuessedPublishedAt().Format("Jan _2 2006"),
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
				fmt.Printf(
					"d---- %*d %s %-*s %s\n",
					longestEpisodes,
					len(r.Episodes()),
					r.GuessedUpdatedAt().Format("Jan _2 2006"),
					longestName,
					r.Name(),
					r.Title(),
				)
			}
		}
	} else {
		rs := onsen.RadioShows()

		cache := make(map[string]adapter.RadioShow)

		for _, r := range rs {
			if r.Name() != "" {
				cache[r.Name()] = r
			}
		}

		for _, arg := range args {
			v, ok := cache[arg]
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
				if e.Manifest() == "" {
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
					"-%s%s%s%s 1 %s %-*s %s\n", playable, bonus, latest, premium,
					e.GuessedPublishedAt().Format("Jan _2 2006"),
					longestName,
					arg+"/"+strconv.FormatUint(uint64(e.EpisodeId()), 10),
					e.Title(),
				)

			}
		}
	}
}
