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

func (r radioShow) GuessedUpdatedAt() time.Time {
	return GuessJstTimeWithNow(r.raw.Updated)
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
