// Package adapter implements onsen.adapter and encapsulates nuxt objects.
package adapter

import (
	"github.com/adios/onsengo/onsen/adapter"
	"github.com/adios/onsengo/onsen/nuxt"
)

// Adapter wraps an instance of nuxt.Nuxt to transform output on the fly.
type nuxtAdapter struct {
	raw *nuxt.Nuxt
}

func (a nuxtAdapter) RadioShows() []adapter.RadioShow {
	all := a.raw.State.Programs.Programs.All

	out := make([]adapter.RadioShow, 0, len(all))
	for i := range all {
		out = append(out, NewRadioShow(&all[i]))
	}
	return out
}

func (a nuxtAdapter) User() *adapter.User {
	if a.raw.State.Signin == nil {
		return nil
	}

	u := NewUser(a.raw.State.Signin)
	return &u
}

func NewAdapter(n *nuxt.Nuxt) adapter.Adapter {
	if n == nil {
		panic("Cannot be nil")
	}
	return nuxtAdapter{n}
}
