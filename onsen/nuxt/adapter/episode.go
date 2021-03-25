package adapter

import (
	"fmt"
	"time"

	"github.com/adios/onsengo/onsen/adapter"
	"github.com/adios/onsengo/onsen/nuxt"
)

type episode struct {
	raw *nuxt.Content
}
type audioEpisode struct{ episode }
type videoEpisode struct{ episode }

func (a audioEpisode) Audio() {}
func (v videoEpisode) Video() {}

func (e episode) EpisodeId() uint {
	return uint(e.raw.Id)
}

func (e episode) RadioShowId() uint {
	return uint(e.raw.ProgramId)
}

func (e episode) Title() string {
	return e.raw.Title
}

func (e episode) Poster() string {
	return e.raw.PosterImageUrl
}

func (e episode) Manifest() string {
	str := e.raw.StreamingUrl
	if str == nil {
		return ""
	}
	return *str
}

func (e episode) PublishedAt() time.Time {
	return GuessJstTimeWithNow(e.raw.DeliveryDate)
}

func (e episode) Guests() []string {
	return e.raw.Guests
}

func (e episode) IsBonus() bool {
	return e.raw.Bonus
}

func (e episode) IsSticky() bool {
	return e.raw.Sticky
}

func (e episode) IsLatest() bool {
	return e.raw.Latest
}

func (e episode) RequiresPremium() bool {
	return e.raw.Premium
}

// episode wraps an instance of nuxt.Content to transform output on the fly.
func NewEpisode(c *nuxt.Content) adapter.Episode {
	if c == nil {
		panic("Cannot be nil")
	}

	var e adapter.Episode

	switch c.MediaType {
	case "sound":
		e = audioEpisode{episode{c}}
	case "movie":
		e = videoEpisode{episode{c}}
	default:
		panic(
			fmt.Sprintf("Content %d: unknown media type\n", c.Id),
		)
	}

	return e
}
