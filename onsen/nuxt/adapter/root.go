// Package adapter encapsulates nuxt.* structs into managed types.
package adapter

import (
	"github.com/adios/onsengo/onsen/adapter"
	"github.com/adios/onsengo/onsen/nuxt"
)

// Root wraps an instance of nuxt.Nuxt to transform output on the fly.
type Root interface {
	RadioShows() []adapter.RadioShow
	// Returns nil if there is no login associated.
	User() *adapter.User
}

type root struct {
	raw *nuxt.Nuxt
}

func (r root) RadioShows() []adapter.RadioShow {
	all := r.raw.State.Programs.Programs.All

	out := make([]adapter.RadioShow, 0, len(all))
	for i := range all {
		out = append(out, NewRadioShow(&all[i]))
	}
	return out
}

func (r root) User() *adapter.User {
	if r.raw.State.Signin == nil {
		return nil
	}

	u := NewUser(r.raw.State.Signin)
	return &u
}

func NewRoot(n *nuxt.Nuxt) Root {
	if n == nil {
		panic("Cannot be nil")
	}
	return root{n}
}
