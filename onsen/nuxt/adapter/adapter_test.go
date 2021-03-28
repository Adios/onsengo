package adapter

import (
	"os"
	"testing"
	"time"

	"github.com/adios/onsengo/onsen/adapter"
	"github.com/adios/onsengo/onsen/nuxt"
	"github.com/stretchr/testify/assert"
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

	anon, _ := nuxt.Parse(string(a))
	paid, _ := nuxt.Parse(string(b))

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

func TestNilNewPerson(t *testing.T) {
	assert.Panics(t, func() { NewPerson(nil) }, "Cannot be nil")
}
func TestPerson(t *testing.T) {
	p := NewPerson(&fts["anon"]["kamisama.performers"].([]nuxt.Performer)[0])

	assert.Equal(t, uint(55), p.PersonId())
	assert.Equal(t, "佐倉綾音", p.Name())
}

func TestNilNewRadioShow(t *testing.T) {
	assert.Panics(t, func() { NewRadioShow(nil) }, "Cannot be nil")
}
func TestAudioRadioShow(t *testing.T) {
	p := fts["anon"]["kamisama"].(nuxt.Program)
	r := NewRadioShow(&p)

	assert.Equal(t, uint(139), r.RadioShowId())
	assert.Equal(t, "kamisama-day", r.Name())
	assert.Equal(t, "神様になったラジオ", r.Title())
	assert.Equal(t, false, r.HasUpdates())

	// Year depends on time.Now()
	at := r.GuessedUpdatedAt()
	assert.Equal(t, time.Month(3), at.Month())
	assert.Equal(t, 19, at.Day())

	h := r.Hosts()
	assert.Equal(t, 2, len(h))
	assert.Equal(t, "佐倉綾音", h[0].Name())
	assert.Equal(t, "花江夏樹", h[1].Name())

	e := r.Episodes()
	assert.Equal(t, 8, len(e))
	assert.Equal(t, "第12回", e[0].Title())
	assert.Equal(t, "第12回 おまけ", e[1].Title())

	_, ok := e[0].(adapter.Audio)
	assert.Equal(t, true, ok)
}
func TestVideoRadioShow(t *testing.T) {
	p := fts["anon"]["fujita"].(nuxt.Program)
	_, ok := NewRadioShow(&p).Episodes()[0].(adapter.Video)

	assert.Equal(t, true, ok)
}

func TestRadioShowWithNoUpdatedTime(t *testing.T) {
	// no updated, but have contents, we take it
	p := fts["anon"]["100man"].(nuxt.Program)
	r := NewRadioShow(&p)

	assert.Equal(t, r.Episodes()[0].GuessedPublishedAt(), r.GuessedUpdatedAt())

	// no updated, no contents (pre-announced show)
	p = fts["anon"]["sorasara"].(nuxt.Program)
	r = NewRadioShow(&p)

	assert.Equal(t, time.Time{}, r.GuessedUpdatedAt())
}

func TestNilNewEpisode(t *testing.T) {
	assert.Panics(t, func() { NewEpisode(nil) }, "Cannot be nil")
}
func TestAudioEpisodeWithGuestNoManifest(t *testing.T) {
	cs := fts["anon"]["kamisama.all"].([]nuxt.Content)
	e := NewEpisode(&cs[6])

	_, ok := e.(adapter.Audio)
	assert.Equal(t, true, ok)

	assert.Equal(t, uint(3114), e.EpisodeId())
	assert.Equal(t, uint(139), e.RadioShowId())
	assert.Equal(t, "第9回", e.Title())

	// Year depends on time.Now()
	at := e.GuessedPublishedAt()
	assert.Equal(t, time.Month(2), at.Month())
	assert.Equal(t, 5, at.Day())

	assert.Equal(
		t,
		"https://d3bzklg4lms4gh.cloudfront.net/program_info/image/default/production"+
			"/66/99/05f3c9402ca36cc3156dd50b7ab9aad298dd/image?v=1602579721",
		e.Poster(),
	)
	assert.Equal(t, "", e.Manifest())
	assert.Equal(
		t,
		[]string{"重松千晴"},
		e.Guests(),
	)

	assert.Equal(t, false, e.IsBonus())
	assert.Equal(t, false, e.IsSticky())
	assert.Equal(t, false, e.IsLatest())
	assert.Equal(t, true, e.RequiresPremium())

}
func TestAudioEpisodeWithManifest(t *testing.T) {
	cs := fts["anon"]["kamisama.all"].([]nuxt.Content)
	e := NewEpisode(&cs[0])

	assert.Equal(t, "HAS_BEEN_SCREENED", e.Manifest())
}

func TestVideoEpisode(t *testing.T) {
	cs := fts["anon"]["fujita.all"].([]nuxt.Content)
	e := NewEpisode(&cs[0])

	_, ok := e.(adapter.Video)
	assert.Equal(t, true, ok)
	assert.Equal(t, "HAS_BEEN_SCREENED", e.Manifest())
}

func TestNologinHasNoUser(t *testing.T) {
	assert.Panics(t, func() { NewUser(fts["anon"]["signin"].(*nuxt.Signin)) }, "Cannot be nil")
}
func TestNewUser(t *testing.T) {
	u := NewUser(fts["paid"]["signin"].(*nuxt.Signin))

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
		u.FollowingRadioShows(),
	)
	assert.Equal(
		t,
		[]uint{3676, 3677},
		u.PlayingEpisodes(),
	)
}

func TestNilNewAdapter(t *testing.T) {
	assert.Panics(t, func() { NewAdapter(nil) }, "Cannot be nil")
}
func TestNewAdapterPaid(t *testing.T) {
	r := NewAdapter(fts["paid"]["nuxt"].(*nuxt.Nuxt))

	assert.Equal(t, 128, len(r.RadioShows()))
	assert.Equal(t, "tsudaken", r.RadioShows()[17].Name())
	assert.NotNil(t, r.User())
}
func TestNewAdapterAnonymous(t *testing.T) {
	r := NewAdapter(fts["anon"]["nuxt"].(*nuxt.Nuxt))

	assert.Equal(t, 128, len(r.RadioShows()))
	assert.Nil(t, r.User())
}

func TestGuessTime(t *testing.T) {
	ref := time.Date(2021, time.Month(3), 24, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, time.Time{}, GuessTime("2020/202/30", ref))

	expected := time.Date(2020, time.Month(3), 25, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, GuessTime("3/25", ref))

	assert.Equal(t, ref.AddDate(0, 0, -1), GuessTime("3/23", ref))
	assert.Equal(t, ref, GuessTime("3/24", ref))
}
