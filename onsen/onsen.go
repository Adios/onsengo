// Package decorator transforms wrapped raw https://onsen.ag/ data (nuxt) for further use.
package decorator

import (
	"regexp"
	"strconv"
	"time"

	"github.com/adios/onsengo/onsen/nuxt"
)

// Transforms nuxt.Nuxt
type Decorator struct {
	Raw *nuxt.Root
}

func (d Decorator) EachRadioShow(fn func(RadioShow)) {
	all := d.Raw.State.Programs.Programs.All

	for i := range all {
		fn(RadioShowFrom(&all[i]))
	}
}

// Always returns a non-nil slice copy.
func (d Decorator) RadioShows() []RadioShow {
	out := make([]RadioShow, 0, len(d.Raw.State.Programs.Programs.All))

	d.EachRadioShow(func(r RadioShow) {
		out = append(out, r)
	})
	return out
}

// Returns an empty User{} if there is no session associated.
func (d Decorator) User() (u User, ok bool) {
	if d.Raw.State.Signin == nil {
		return User{}, false
	}
	return UserFrom(d.Raw.State.Signin), true
}

func DecoratorFrom(n *nuxt.Root) Decorator {
	if n == nil {
		panic("Cannot be nil")
	}
	return Decorator{n}
}

// Transforms nuxt.Program
type RadioShow struct {
	Raw *nuxt.Program
}

func (r RadioShow) RadioShowId() RadioShowId {
	return RadioShowId(r.Raw.Id)
}

func (r RadioShow) Name() string {
	return r.Raw.DirectoryName
}

func (r RadioShow) Title() string {
	return r.Raw.Title
}

func (r RadioShow) HasBeenUpdated() bool {
	return r.Raw.New
}

// Returns the time in which its year component (YYYY) is never after now year.
// It's doing so because onsen.ag gives the value without a year component, i.e. only "MM/DD".
// We have to guess the YYYY by ourself in order to make this field useful.
//
// If the raw value isn't in MM/DD format, an time.Time{} will be returned.
//
// The timezone associated is UTC+9 for all successful returns.
//
// BUG(adios): It's possible we returns a time with wrong YYYY value.
func (r RadioShow) JstUpdatedAt() (res time.Time, ok bool) {
	t := r.Raw.Updated

	// Try using first (latest) episode's MM/DD if current show has no MM/DD
	if t == nil {
		if cs := r.Raw.Contents; len(cs) == 0 || cs[0].DeliveryDate == "" {
			return time.Time{}, false
		}
		t = &r.Raw.Contents[0].DeliveryDate
	}
	return GuessJstTimeWithNow(*t)
}

func (r RadioShow) EachHost(fn func(host Person)) {
	ps := r.Raw.Performers

	for i := range ps {
		fn(PersonFrom(&ps[i]))
	}
}

// Always returns a non-nil slice copy.
func (r RadioShow) Hosts() []Person {
	out := make([]Person, 0, len(r.Raw.Performers))

	r.EachHost(func(p Person) {
		out = append(out, p)
	})
	return out
}

func (r RadioShow) EachEpisode(fn func(Episode)) {
	cs := r.Raw.Contents

	for i := range cs {
		fn(EpisodeFrom(&cs[i]))
	}
}

// Always returns a non-nil slice copy.
func (r RadioShow) Episodes() []Episode {
	out := make([]Episode, 0, len(r.Raw.Contents))

	r.EachEpisode(func(e Episode) {
		out = append(out, e)
	})
	return out
}

func RadioShowFrom(p *nuxt.Program) RadioShow {
	if p == nil {
		panic("Cannot be nil")
	}
	return RadioShow{p}
}

// Transforms nuxt.Content
type Episode struct {
	Raw *nuxt.Content
}

func (e Episode) EpisodeId() EpisodeId {
	return EpisodeId(e.Raw.Id)
}

func (e Episode) RadioShowId() RadioShowId {
	return RadioShowId(e.Raw.ProgramId)
}

func (e Episode) Title() string {
	return e.Raw.Title
}

// The URL of episode's poster image.
func (e Episode) Poster() (url string) {
	return e.Raw.PosterImageUrl
}

// The URL of episode's m3u8 manifest. An empty string means the resource is not accessible with current session.
func (e Episode) Manifest() (url string, ok bool) {
	str := e.Raw.StreamingUrl

	if str == nil {
		return "", false
	}
	return *str, true
}

// Returns the time in which its year component (YYYY) is never after now year.
// It's doing so because onsen.ag gives the value without a year component, i.e. only "MM/DD".
// We have to guess the YYYY by ourself in order to make this field useful.
//
// If the raw value isn't in MM/DD format, an time.Time{} will be returned.
//
// The timezone associated is UTC+9 for all successful returns.
//
// BUG(adios): It's possible we returns a time with wrong YYYY value.
func (e Episode) JstUpdatedAt() (res time.Time, ok bool) {
	return GuessJstTimeWithNow(e.Raw.DeliveryDate)
}

func (e Episode) EachGuest(fn func(name string)) {
	for _, g := range e.Raw.Guests {
		fn(g)
	}
}

// Always returns a non-nil slice copy.
func (e Episode) Guests() (names []string) {
	out := make([]string, len(e.Raw.Guests))

	// safe to copy, e.raw.Guests will never be a nil slice.
	copy(out, e.Raw.Guests)
	return out
}

func (e Episode) IsBonus() bool {
	return e.Raw.Bonus
}

func (e Episode) IsSticky() bool {
	return e.Raw.Sticky
}

func (e Episode) IsLatest() bool {
	return e.Raw.Latest
}

func (e Episode) RequiresPremium() bool {
	return e.Raw.Premium
}

func (e Episode) HasVideoStream() bool {
	return e.Raw.Movie
}

func EpisodeFrom(c *nuxt.Content) Episode {
	if c == nil {
		panic("Cannot be nil")
	}
	return Episode{c}
}

// Transforms nuxt.Signin
type User struct {
	Raw *nuxt.Signin
}

func (u User) Email() string {
	return u.Raw.Email
}

func (u User) UserId() UserId {
	return UserId(u.Raw.Id)
}

func (u User) EachFollowingPerson(fn func(id PersonId)) {
	for _, id := range u.Raw.FavoritePerformerIds {
		fn(PersonId(id))
	}
}

// Always returns a non-nil slice copy.
func (u User) FollowingPeople() []PersonId {
	out := make([]PersonId, 0, len(u.Raw.FavoritePerformerIds))

	u.EachFollowingPerson(func(id PersonId) {
		out = append(out, id)
	})
	return out
}

func (u User) EachFollowingShow(fn func(id RadioShowId)) {
	for _, id := range u.Raw.FavoriteProgramIds {
		fn(RadioShowId(id))
	}
}

// Always returns a non-nil slice copy.
func (u User) FollowingShows() []RadioShowId {
	out := make([]RadioShowId, 0, len(u.Raw.FavoriteProgramIds))

	u.EachFollowingShow(func(id RadioShowId) {
		out = append(out, id)
	})
	return out
}

func (u User) EachPlaylistEpisode(fn func(id EpisodeId)) {
	for _, id := range u.Raw.PlaylistedContentIds {
		fn(EpisodeId(id))
	}
}

// Always returns a non-nil slice copy.
func (u User) PlaylistEpisodes() []EpisodeId {
	out := make([]EpisodeId, 0, len(u.Raw.PlaylistedContentIds))

	u.EachPlaylistEpisode(func(id EpisodeId) {
		out = append(out, id)
	})
	return out
}

func UserFrom(s *nuxt.Signin) User {
	if s == nil {
		panic("Cannot be nil")
	}
	return User{s}
}

// Transforms nuxt.Performer
type Person struct {
	Raw *nuxt.Performer
}

func (p Person) PersonId() PersonId {
	return PersonId(p.Raw.Id)
}

func (p Person) Name() string {
	return p.Raw.Name
}

func PersonFrom(p *nuxt.Performer) Person {
	if p == nil {
		panic("Cannot be nil")
	}
	return Person{p}
}

type EpisodeId uint

func (id EpisodeId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type PersonId uint

func (id PersonId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type RadioShowId uint

func (id RadioShowId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type UserId string

// Given a date string with no YYYY component (MM/DD) and a referenced time (usually now),
// we find a most recent year (YYYY) such that YYYY/MM/DD won't go over the referenced time.
func GuessTime(guess string, ref time.Time) (res time.Time, ok bool) {
	re := regexp.MustCompile("^([0-9]{1,2})/([0-9]{1,2})$")
	m := re.FindStringSubmatch(guess)

	if m == nil {
		return time.Time{}, false
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
		ref.Year(),
		time.Month(guessMonth),
		guessDay, 0, 0, 0, 0,
		ref.Location(),
	)

	if attemptTime.After(ref) {
		return attemptTime.AddDate(-1, 0, 0), true
	} else {
		return attemptTime, true
	}
}

// Set a zone of UTC+9 on top of GuessTime().
func GuessJstTimeWithNow(guess string) (res time.Time, ok bool) {
	loc := time.FixedZone("UTC+9", 9*60*60)
	now := time.Now().In(loc)

	return GuessTime(guess, now)
}
package onsen

import (
	"regexp"
)

// Returns a byte slice to the capture of first appeared NUXT pattern:
//   <script>window.__NUXT__=([^<]+);</script>
func FindNuxtExpression(html []byte) (expr []byte, ok bool) {
	re := regexp.MustCompile("<script>window.__NUXT__=([^<]+);</script>")

	m := re.FindSubmatch(html)
	if m == nil {
		return nil, false
	}
	return m[1], true
}
// Package expression is a wrapper of github.com/dop251/goja to run Javascript code (expressions) for deobfuscation.
package expression

import (
	"fmt"

	"github.com/dop251/goja"
)

type Expression struct {
	js string
	vm *goja.Runtime
}

// Run the given JavaScript code which can produces a *value*, i.e. expressions.
// Returns a string of the value's JSON representation and any JS error encountered.
// Note that "undefined" is also considered as an error.
func (e *Expression) Stringify() (json string, err error) {
	torun := fmt.Sprintf("JSON.stringify(%s)", string(e.js))

	res, err := e.getVm().RunString(torun)
	if err != nil {
		return "", err
	}

	out := res.Export()
	if out == nil {
		return "", fmt.Errorf("Got nothing after running. Possibly the js returned an undefined.\n")
	}

	return out.(string), nil
}

func (e *Expression) getVm() *goja.Runtime {
	if e.vm == nil {
		e.vm = goja.New()
	}
	return e.vm
}

func From(js string) *Expression {
	return &Expression{
		js: js,
	}
}
