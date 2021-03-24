package adapter

import (
	"github.com/adios/onsengo/onsen/nuxt"
)

// User wraps an instance of nuxt.Signin to transform output on the fly.
type User interface {
	Email() string
	// Raw value is a string of digits, possibly int type. But we don't know, leave it as is.
	UserId() string
	// A slice of Person.PersonId()
	FollowingPeople() []uint
	// A slice of RadioShow.RadioShowId()
	FollowingRadioShows() []uint
	// A slice of Episode.EpisodeId()
	PlayingEpisodes() []uint
}

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

func NewUser(s *nuxt.Signin) User {
	if s == nil {
		panic("Cannot be nil")
	}
	return user{s}
}
