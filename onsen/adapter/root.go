// Package adapter encapsulates nuxt.* structs into managed types.
package adapter

import (
	"github.com/adios/onsengo/onsen/nuxt"
)

// Root wraps an instance of nuxt.Nuxt to transform output on the fly.
type Root interface {
	RadioShows() []RadioShow
	// Returns nil if there is no login associated.
	User() *User
}

type root struct {
	raw *nuxt.Nuxt
}

func (r root) RadioShows() []RadioShow {
	all := r.raw.State.Programs.Programs.All

	out := make([]RadioShow, 0, len(all))
	for i := range all {
		out = append(out, NewRadioShow(&all[i]))
	}
	return out
}

func (r root) User() *User {
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
