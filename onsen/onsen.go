// Package onsen embeds a Nuxt decorator acting as a front-end to clients.
//
// SIDE EFFECT:
//
// The following function and methods have side effects:
//
//    Radio.JstUpdatedAt()
//    Episode.JstUpdatedAt()
//    GuessJstTimeWithNow()
//
// Their outputs depend on time.Now().
//
// Set a fixed time if need to test thier output values:
//
//     // any fixed date is OK as long as it fits the data getting test.
//     onsen.SetRefDate("2021-03-21")
//
// In raw json, the upload date of all radio shows is in a string of MM/DD format,
// in order to build a complete timestamp for the date, they must do a guess to find a possible YYYY.
// By calling SetRefDate(), it sets a fixed date to guess intead of time.Now().
package onsen

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/dop251/goja"

	"github.com/adios/onsengo/onsen/nuxt"
)

// Set this to a fixed time to test GuessJstTimeWithNow() and JstUpdatedAt().
var guessRefTime = time.Now()

func SetRefDate(date string) {
	tm, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	guessRefTime = tm
}

type index map[interface{}]Radio

type Onsen struct {
	Nuxt
	cache index
}

func (o *Onsen) Radio(id interface{}) (r Radio, ok bool) {
	r, ok = o.Index()[id]
	return
}

func (o *Onsen) Index() index {
	if o.cache == nil {
		c := make(index)
		o.EachRadio(func(r Radio) {
			c[r.Id()] = r
			c[r.Name()] = r
		})
		o.cache = c
	}
	return o.cache
}

func Create(html string) (*Onsen, error) {
	expr, ok := FindNuxtExpression(html)
	if !ok {
		return nil, fmt.Errorf("Create: NUXT pattern not matched")
	}

	str, err := StringifyExpression(expr)
	if err != nil {
		return nil, err
	}

	n, err := nuxt.Create(str)
	if err != nil {
		return nil, err
	}

	return &Onsen{
		Nuxt: Nuxt{n},
	}, nil
}

// Transforms nuxt.Nuxt.
type Nuxt struct {
	Raw *nuxt.Nuxt
}

func (n Nuxt) EachRadio(fn func(Radio)) {
	rs := n.programs()

	for i := range rs {
		fn(Radio{&rs[i]})
	}
}

// Returns a new copy of non-nil slice.
func (n Nuxt) Radios() []Radio {
	rs := n.programs()
	out := make([]Radio, len(rs))

	for i := range rs {
		out[i] = Radio{&rs[i]}
	}
	return out
}

// Returns an empty User{} if there is no session associated.
func (n Nuxt) User() (u User, ok bool) {
	if n.Raw.State.Signin == nil {
		return User{}, false
	}
	return User{n.Raw.State.Signin}, true
}

func (n Nuxt) programs() []nuxt.Program {
	return n.Raw.State.Programs.Programs.All
}

// Transforms nuxt.Program.
type Radio struct {
	Raw *nuxt.Program
}

func (r Radio) Id() int {
	return r.Raw.Id
}

func (r Radio) Name() string {
	return r.Raw.DirectoryName
}

func (r Radio) Title() string {
	return r.Raw.Title
}

func (r Radio) HasBeenUpdated() bool {
	return r.Raw.New
}

// SIDE EFFECT: this method has side effect, set a fixed year by SetRefDate() when testing its value.
//
// If a show is updated on "3/19", the method returns a time with its date set on either "2021/03/19" or "2020/03/19",
// depends on time.Now(). Since there is no year component given in the raw output from onsen.ag.
//
// If the raw value isn't in MM/DD format, an time.Time{} will be returned.
//
// The timezone associated is UTC+9 for all successful returns.
//
// BUG(adios): It's possible we returns a time with wrong YYYY value.
func (r Radio) JstUpdatedAt() (res time.Time, ok bool) {
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

// Returns a new copy of non-nil slice.
func (r Radio) Hosts() []Person {
	out := make([]Person, len(r.Raw.Performers))
	for i := range r.Raw.Performers {
		out[i] = Person{&r.Raw.Performers[i]}
	}
	return out
}

// Returns a new copy of non-nil slice.
func (r Radio) Episodes() []Episode {
	out := make([]Episode, len(r.Raw.Contents))
	for i := range r.Raw.Contents {
		out[i] = Episode{&r.Raw.Contents[i]}
	}
	return out
}

// Transforms nuxt.Content.
type Episode struct {
	Raw *nuxt.Content
}

func (e Episode) Id() int {
	return e.Raw.Id
}

func (e Episode) RadioId() int {
	return e.Raw.ProgramId
}

func (e Episode) Title() string {
	return e.Raw.Title
}

// The URL to episode's poster image.
func (e Episode) Poster() (url string) {
	return e.Raw.PosterImageUrl
}

// The URL to episode's m3u8 manifest. An empty string means the resource is not accessible with current session.
func (e Episode) Manifest() (url string, ok bool) {
	str := e.Raw.StreamingUrl

	if str == nil {
		return "", false
	}
	return *str, true
}

// SIDE EFFECT: this method has side effect, set a fixed year by SetRefDate() when testing its value.
//
// If a show is updated on "3/19", the method returns a time with its date set on either "2021/03/19" or "2020/03/19",
// depends on time.Now(). Since there is no year component given in the raw output from onsen.ag.
//
// If the raw value isn't in MM/DD format, an time.Time{} will be returned.
//
// The timezone associated is UTC+9 for all successful returns.
//
// BUG(adios): It's possible we returns a time with wrong YYYY value.
func (e Episode) JstUpdatedAt() (res time.Time, ok bool) {
	return GuessJstTimeWithNow(e.Raw.DeliveryDate)
}

// Returns a new copy of non-nil slice.
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

// Transforms nuxt.Signin.
type User struct {
	Raw *nuxt.Signin
}

func (u User) Email() string {
	return u.Raw.Email
}

func (u User) Id() string {
	return u.Raw.Id
}

// Returns a new copy of non-nil slice.
func (u User) FollowingPeople() []int {
	out := make([]int, len(u.Raw.FavoritePerformerIds))
	copy(out, u.Raw.FavoritePerformerIds)
	return out
}

// Returns a new copy of non-nil slice.
func (u User) FollowingRadios() []int {
	out := make([]int, len(u.Raw.FavoriteProgramIds))
	copy(out, u.Raw.FavoriteProgramIds)
	return out
}

// Returns a new copy of non-nil slice.
func (u User) PlaylistEpisodes() []int {
	out := make([]int, len(u.Raw.PlaylistedContentIds))
	copy(out, u.Raw.PlaylistedContentIds)
	return out
}

// Transforms nuxt.Performer.
type Person struct {
	Raw *nuxt.Performer
}

func (p Person) Id() int {
	return p.Raw.Id
}

func (p Person) Name() string {
	return p.Raw.Name
}

// Run the given JavaScript code for deobfuscation.
// The code must produce a *value*, i.e. expressions.
// Returns a string of the value's JSON representation and any JS error encountered.
// Note that "undefined" is also considered as an error.
func StringifyExpression(expr string) (string, error) {
	js := fmt.Sprintf("JSON.stringify(%s)", expr)

	res, err := goja.New().RunString(js)
	if err != nil {
		return "", err
	}

	out := res.Export()
	if out == nil {
		return "", fmt.Errorf("StringifyExpression: possibly js returned an undefined")
	}
	return out.(string), nil
}

// Returns a string to the capture of first appeared NUXT pattern:
//   <script>window.__NUXT__=([^<]+);</script>
func FindNuxtExpression(html string) (expr string, ok bool) {
	re := regexp.MustCompile("<script>window.__NUXT__=([^<]+);</script>")

	m := re.FindStringSubmatch(html)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// Given a date string with no YYYY component (MM/DD) and a referenced time,
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

// SIDE EFFECT: this method has side effect, set a fixed year by SetRefDate() when testing its value.
//
// Set UTC+9 fixed time zone on top of GuessTime().
func GuessJstTimeWithNow(guess string) (res time.Time, ok bool) {
	loc := time.FixedZone("UTC+9", 9*60*60)
	ref := guessRefTime.In(loc)

	return GuessTime(guess, ref)
}
