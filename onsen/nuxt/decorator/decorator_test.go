package decorator

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/adios/onsengo/onsen/nuxt"
)

var fts map[string]map[string]interface{}

func setupFixtures() {
	a, err := os.ReadFile("../testdata/fixture_nologin_screened.json")
	if err != nil {
		panic(err)
	}
	b, err := os.ReadFile("../testdata/fixture_paid_screened.json")
	if err != nil {
		panic(err)
	}

	anon, _ := nuxt.FromReader(bytes.NewReader(a))
	paid, _ := nuxt.FromReader(bytes.NewReader(b))

	anonAll := anon.State.Programs.Programs.All
	paidAll := paid.State.Programs.Programs.All

	fts = map[string]map[string]interface{}{
		"anon": {
			"nuxt":                anon,
			"programs":            anonAll,
			"kamisama":            anonAll[34],
			"kamisama.all":        anonAll[34].Contents,
			"kamisama.performers": anonAll[34].Performers,
			"fujita":              anonAll[17],
			"fujita.all":          anonAll[17].Contents,
			"100man":              anonAll[0],
			"sorasara":            anonAll[3],
			"signin":              anon.State.Signin,
		},
		"paid": {
			"nuxt":     paid,
			"programs": paidAll,
			"signin":   paid.State.Signin,
		},
	}

}

func TestMain(m *testing.M) {
	setupFixtures()
	ret := m.Run()
	os.Exit(ret)
}

func TestPersonFromNil(t *testing.T) {
	assert.Panics(t, func() { PersonFrom(nil) }, "Cannot be nil")
}
func TestPerson(t *testing.T) {
	p := PersonFrom(&fts["anon"]["kamisama.performers"].([]nuxt.Performer)[0])

	assert.Equal(t, uint(55), p.PersonId())
	assert.Equal(t, "佐倉綾音", p.Name())
}

func TestRadioShowFromNil(t *testing.T) {
	assert.Panics(t, func() { RadioShowFrom(nil) }, "Cannot be nil")
}
func TestRadioShow(t *testing.T) {
	p := fts["anon"]["kamisama"].(nuxt.Program)
	r := RadioShowFrom(&p)

	assert.Equal(t, uint(139), r.RadioShowId())
	assert.Equal(t, "kamisama-day", r.Name())
	assert.Equal(t, "神様になったラジオ", r.Title())
	assert.False(t, r.HasBeenUpdated())

	// Skip check year as it's dependent on time.Now().
	at, ok := r.JstUpdatedAt()
	assert.True(t, ok)
	assert.Equal(t, time.Month(3), at.Month())
	assert.Equal(t, 19, at.Day())

	name, offset := at.Zone()
	assert.Equal(t, "UTC+9", name)
	assert.Equal(t, 9*60*60, offset)

	h := r.Hosts()
	assert.Equal(t, 2, len(h))
	assert.Equal(t, "佐倉綾音", h[0].Name())
	assert.Equal(t, "花江夏樹", h[1].Name())

	e := r.Episodes()
	assert.Equal(t, 8, len(e))
	assert.Equal(t, "第12回", e[0].Title())
	assert.Equal(t, "第12回 おまけ", e[1].Title())
}

func TestRadioShowWithNoUpdatedTime(t *testing.T) {
	// no updated, but have contents, we take it
	p := fts["anon"]["100man"].(nuxt.Program)
	r := RadioShowFrom(&p)

	showTime, ok := r.JstUpdatedAt()
	assert.True(t, ok)

	epTime, ok := r.Episodes()[0].JstUpdatedAt()
	assert.True(t, ok)
	assert.Equal(t, epTime, showTime)

	// no updated, no contents (pre-announced show)
	p = fts["anon"]["sorasara"].(nuxt.Program)
	r = RadioShowFrom(&p)

	showTime, ok = r.JstUpdatedAt()
	assert.False(t, ok)
	assert.Equal(t, time.Time{}, showTime)
}

func TestEpisodeFromNil(t *testing.T) {
	assert.Panics(t, func() { EpisodeFrom(nil) }, "Cannot be nil")
}
func TestEpisodeWithAudioGuestsButManifest(t *testing.T) {
	cs := fts["anon"]["kamisama.all"].([]nuxt.Content)
	e := EpisodeFrom(&cs[6])

	assert.Equal(t, uint(3114), e.EpisodeId())
	assert.Equal(t, uint(139), e.RadioShowId())
	assert.Equal(t, "第9回", e.Title())

	// Skip check year as it's dependent on time.Now().
	at, ok := e.JstUpdatedAt()
	assert.True(t, ok)
	assert.Equal(t, time.Month(2), at.Month())
	assert.Equal(t, 5, at.Day())

	name, offset := at.Zone()
	assert.Equal(t, "UTC+9", name)
	assert.Equal(t, 9*60*60, offset)

	assert.Equal(
		t,
		"https://d3bzklg4lms4gh.cloudfront.net/program_info/image/default/production"+
			"/66/99/05f3c9402ca36cc3156dd50b7ab9aad298dd/image?v=1602579721",
		e.Poster(),
	)

	m, ok := e.Manifest()
	assert.False(t, ok)
	assert.Equal(t, "", m)

	assert.Equal(
		t,
		[]string{"重松千晴"},
		e.Guests(),
	)

	assert.False(t, e.IsBonus())
	assert.False(t, e.IsSticky())
	assert.False(t, e.IsLatest())
	assert.True(t, e.RequiresPremium())
	assert.False(t, e.HasVideoStream())
}
func TestEpisodeWithManifest(t *testing.T) {
	cs := fts["anon"]["kamisama.all"].([]nuxt.Content)
	e := EpisodeFrom(&cs[0])

	m, ok := e.Manifest()
	assert.True(t, ok)
	assert.Equal(t, "HAS_BEEN_SCREENED", m)
}

func TestEpisodeWithVideo(t *testing.T) {
	cs := fts["anon"]["fujita.all"].([]nuxt.Content)
	e := EpisodeFrom(&cs[0])

	assert.True(t, e.HasVideoStream())
}

func TestUserFromNil(t *testing.T) {
	assert.Panics(t, func() { UserFrom(fts["anon"]["signin"].(*nuxt.Signin)) }, "Cannot be nil")
}
func TestNewUser(t *testing.T) {
	u := UserFrom(fts["paid"]["signin"].(*nuxt.Signin))

	assert.Equal(t, "hello@world", u.Email())
	assert.Equal(t, "0", u.UserId())
	assert.Equal(
		t,
		[]uint{
			1377, 1044, 946, 889, 726, 645, 641,
			590, 559, 429, 421, 284, 211, 136, 114,
			113, 105, 77, 66, 55, 396, 29,
		},
		u.FollowingPeople(),
	)
	assert.Equal(
		t,
		[]uint{
			4, 10, 16, 17, 18, 29, 47, 54, 56, 65, 76,
			77, 88, 89, 93, 118, 131, 136, 139, 149, 156, 159,
		},
		u.FollowingShows(),
	)
	assert.Equal(
		t,
		[]uint{3676, 3677},
		u.PlaylistEpisodes(),
	)
}

func TestAdapterFromNil(t *testing.T) {
	assert.Panics(t, func() { DecoratorFrom(nil) }, "Cannot be nil")
}
func TestAdapterFromPaid(t *testing.T) {
	r := DecoratorFrom(fts["paid"]["nuxt"].(*nuxt.Root))

	assert.Equal(t, 128, len(r.RadioShows()))
	assert.Equal(t, "tsudaken", r.RadioShows()[17].Name())

	u, ok := r.User()
	assert.True(t, ok)
	assert.NotNil(t, u)
}
func TestAdapterFromAnonymous(t *testing.T) {
	r := DecoratorFrom(fts["anon"]["nuxt"].(*nuxt.Root))

	assert.Equal(t, 128, len(r.RadioShows()))

	u, ok := r.User()
	assert.False(t, ok)
	assert.Equal(t, User{}, u)
}

func TestGuessTime(t *testing.T) {
	ref := time.Date(2021, time.Month(3), 24, 0, 0, 0, 0, time.UTC)

	g, ok := GuessTime("2020/202/30", ref)
	assert.False(t, ok)
	assert.Equal(t, time.Time{}, g)

	expected := time.Date(2020, time.Month(3), 25, 0, 0, 0, 0, time.UTC)
	g, ok = GuessTime("3/25", ref)
	assert.True(t, ok)
	assert.Equal(t, expected, g)

	g, ok = GuessTime("3/23", ref)
	assert.True(t, ok)
	assert.Equal(t, ref.AddDate(0, 0, -1), g)

	g, ok = GuessTime("3/24", ref)
	assert.True(t, ok)
	assert.Equal(t, ref, g)
}
