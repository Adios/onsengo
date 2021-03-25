package adapter

import (
	"github.com/adios/onsengo/onsen/adapter"
	"github.com/adios/onsengo/onsen/nuxt"
)

type user struct {
	raw *nuxt.Signin
}

func (u user) Email() string {
	return u.raw.Email
}

func (u user) UserId() string {
	return u.raw.Id
}

func (u user) FollowingPeople() []uint {
	ps := u.raw.FavoritePerformerIds

	out := make([]uint, 0, len(ps))
	for _, pid := range ps {
		out = append(out, uint(pid))
	}
	return out
}

func (u user) FollowingRadioShows() []uint {
	ps := u.raw.FavoriteProgramIds

	out := make([]uint, 0, len(ps))
	for _, pid := range ps {
		out = append(out, uint(pid))
	}
	return out
}

func (u user) PlayingEpisodes() []uint {
	cs := u.raw.PlaylistedContentIds

	out := make([]uint, 0, len(cs))
	for _, cid := range cs {
		out = append(out, uint(cid))
	}
	return out
}

// User wraps an instance of nuxt.Signin to transform output on the fly.
func NewUser(s *nuxt.Signin) adapter.User {
	if s == nil {
		panic("Cannot be nil")
	}
	return user{s}
}
