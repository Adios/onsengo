package adapter

import (
	"time"

	"github.com/adios/onsengo/onsen/adapter"
	"github.com/adios/onsengo/onsen/nuxt"
)

type radioShow struct {
	raw *nuxt.Program
}

func (r radioShow) RadioShowId() uint {
	return uint(r.raw.Id)
}

func (r radioShow) Name() string {
	return r.raw.DirectoryName
}

func (r radioShow) Title() string {
	return r.raw.Title
}

func (r radioShow) HasUpdates() bool {
	return r.raw.New
}

// Returns a best-effor time that is guessed based on time.Now().
// Since there is no YYYY recorded in onsen's raw data. (MM/DD only)
// An empty time.Time{} means there is an invalid date pattern or just not having a time.
func (r radioShow) GuessedUpdatedAt() time.Time {
	t := r.raw.Updated
	if t == nil {
		// Null updated can be found if:
		// * the radio is just announced, not having content yet, or
		// * it got re-announced and they didn't set an updated time
		//
		// Manually set to the latest one if it has contents
		if cs := r.raw.Contents; len(cs) == 0 || cs[0].DeliveryDate == "" {
			return time.Time{}
		}
		t = &r.raw.Contents[0].DeliveryDate
	}
	return GuessJstTimeWithNow(*t)
}

func (r radioShow) Hosts() []adapter.Person {
	ps := r.raw.Performers

	out := make([]adapter.Person, 0, len(ps))
	for i := range ps {
		out = append(out, NewPerson(&ps[i]))
	}

	return out
}

func (r radioShow) Episodes() []adapter.Episode {
	cs := r.raw.Contents

	out := make([]adapter.Episode, 0, len(cs))
	for i := range cs {
		out = append(out, NewEpisode(&cs[i]))
	}

	return out
}

// RadioShow wraps an instance of nuxt.Program to transform output on the fly.
func NewRadioShow(p *nuxt.Program) adapter.RadioShow {
	if p == nil {
		panic("Cannot be nil")
	}
	return radioShow{p}
}
