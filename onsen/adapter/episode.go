package adapter

import (
	"fmt"
	"time"

	"github.com/adios/onsengo/onsen/nuxt"
)

// Episode wraps an instance of nuxt.Content to transform output on the fly.
type Episode interface {
	EpisodeId() uint
	RadioShowId() uint
	Title() string

	// Returns the url to the episode's poster image.
	Poster() string

	// Returns the url to the episode's manifest (m3u8).
	// An empty string ("") means the resource is not accessible from current user identity.
	Manifest() string

	// Original value is an incomplete date string in MM/DD format, i.e. no year specified.
	// This method try to guess a best effor year based on time.Now().
	// An empty time.Time{} means there is an invalid date pattern.
	PublishedAt() time.Time

	// Always returns a slice, an empty slice means there are no guests.
	//  [ "name1", "name2" ]
	Guests() []string

	IsBonus() bool
	IsSticky() bool
	IsLatest() bool
	RequiresPremium() bool
}

// Audio is attached to the Episodes that deliver sound-only contents. You can check an Episode with type assertion:
//  _, ok = episode.(Audio)
// There are no public methods for now.
type Audio interface {
	audio()
}

// Video is attached to the Episodes that deliver video contents. You can check an Episode with type assertion:
//  _, ok = episode.(Video)
// There are no public methods for now.
type Video interface {
	video()
}

type episode struct {
	raw *nuxt.Content
}
type audioEpisode struct{ episode }
type videoEpisode struct{ episode }

func (a audioEpisode) audio() {}
func (v videoEpisode) video() {}

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

func NewEpisode(c *nuxt.Content) Episode {
	if c == nil {
		panic("Cannot be nil")
	}

	var e Episode

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
