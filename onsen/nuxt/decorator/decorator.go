// Package adapter implements onsen.adapter and encapsulates nuxt objects.
package decorator

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

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

func (a nuxtAdapter) User() adapter.User {
	if a.raw.State.Signin == nil {
		return nil
	}

	u := NewUser(a.raw.State.Signin)
	return u
}

func NewAdapter(n *nuxt.Nuxt) adapter.Adapter {
	if n == nil {
		panic("Cannot be nil")
	}
	return nuxtAdapter{n}
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

// User wraps an instance of nuxt.Signin to transform output on the fly.
func NewUser(s *nuxt.Signin) adapter.User {
	if s == nil {
		panic("Cannot be nil")
	}
	return user{s}
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

type person struct {
	raw *nuxt.Performer
}

func (p person) PersonId() uint {
	return uint(p.raw.Id)
}

func (p person) Name() string {
	return p.raw.Name
}

// Person wraps an instance of nuxt.Performer to transform output on the fly.
func NewPerson(p *nuxt.Performer) adapter.Person {
	if p == nil {
		panic("Cannot be nil")
	}
	return person{p}
}

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

func (e episode) GuessedPublishedAt() time.Time {
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

// Given an incomplete date string in the format "MM/DD", and with a referenced date.
// The function finds a nearest year such that "YYYY/MM/DD" doesn't go over the referenced date.
func GuessTime(guess string, reference time.Time) time.Time {
	re := regexp.MustCompile("^([0-9]{1,2})/([0-9]{1,2})$")
	m := re.FindStringSubmatch(guess)

	if m == nil {
		return time.Time{}
	}

	guessMonth, err := strconv.Atoi(m[1])
	if err != nil {
		panic(err)
	}
	guessDay, err := strconv.Atoi(m[2])
	if err != nil {
		panic(err)
	}

	attemptTime := time.Date(
		reference.Year(),
		time.Month(guessMonth),
		guessDay, 0, 0, 0, 0,
		reference.Location(),
	)

	if attemptTime.After(reference) {
		return attemptTime.AddDate(-1, 0, 0)
	} else {
		return attemptTime
	}
}

// Based on the GuessTime() function, here we set a timezone of UTC+9 and use time.Now() as a referenced date.
func GuessJstTimeWithNow(guess string) time.Time {
	loc := time.FixedZone("UTC+9", 9*60*60)
	now := time.Now().In(loc)

	return GuessTime(guess, now)
}
