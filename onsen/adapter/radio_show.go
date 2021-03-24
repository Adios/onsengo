package adapter

import (
	"time"

	"github.com/adios/onsengo/onsen/nuxt"
)

// RadioShow wraps an instance of nuxt.Program to transform output on the fly.
type RadioShow interface {
	RadioShowId() uint
	Name() string
	Title() string
	HasUpdates() bool
	UpdatedAt() time.Time
	Hosts() []Person

	// Returns a slice of Episode instances which may either be an AudioEpisode or a VideoEpisode.
	Episodes() []Episode
}

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

func (r radioShow) UpdatedAt() time.Time {
	return GuessJstTimeWithNow(r.raw.Updated)
}

func (r radioShow) Hosts() []Person {
	ps := r.raw.Performers

	out := make([]Person, 0, len(ps))
	for i := range ps {
		out = append(out, NewPerson(&ps[i]))
	}

	return out
}

func (r radioShow) Episodes() []Episode {
	cs := r.raw.Contents

	out := make([]Episode, 0, len(cs))
	for i := range cs {
		out = append(out, NewEpisode(&cs[i]))
	}

	return out
}

func NewRadioShow(p *nuxt.Program) RadioShow {
	if p == nil {
		panic("Cannot be nil")
	}
	return radioShow{p}
}
