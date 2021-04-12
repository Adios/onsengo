package onsen

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/adios/onsengo/onsen/nuxt"
)

func TestStringifyExpression(t *testing.T) {
	assert := assert.New(t)

	{
		s, err := StringifyExpression("")
		assert.Empty(s)
		assert.EqualError(err, "StringifyExpression: possibly js returned an undefined")
	}
	{
		s, err := StringifyExpression(";")
		assert.Empty(s)
		assert.Error(err)
		assert.Contains(err.Error(), "Unexpected token")
	}
}

func TestFindNuxtExpression(t *testing.T) {
	tests := map[string]struct {
		in  string
		out string
		ok  bool
	}{
		"empty str": {"", "", false},
		"invalid pattern": {
			"...<script>window.__NUXT__=one</script><script>window.__NUXT__=two</script>...", "", false,
		},
		"ok pattern": {
			"...<script>window.__NUXT__=one</script><script>window.__NUXT__=two;</script>...", "two", true,
		},
	}
	for _, test := range tests {
		out, ok := FindNuxtExpression(test.in)
		assert.Equal(t, test.out, out)
		assert.Equal(t, test.ok, ok)
	}
}

func TestGuessTime(t *testing.T) {
	var (
		pt = func(str string) time.Time {
			out, _ := time.Parse("2006-01-02", str)
			return out
		}
		ref = pt("2021-03-24")
	)

	tests := []struct {
		in  string
		out time.Time
		ok  bool
	}{
		{"2020/202/30", time.Time{}, false},
		{"3/25", pt("2020-03-25"), true},
		{"3/23", pt("2021-03-23"), true},
		{"3/24", pt("2021-03-24"), true},
	}
	for _, test := range tests {
		out, ok := GuessTime(test.in, ref)
		assert.Equal(t, test.out, out)
		assert.Equal(t, test.ok, ok)
	}
}

func TestNuxtWithAnonymousUser(t *testing.T) {
	var (
		assert = assert.New(t)
		f, _   = os.ReadFile("testdata/fixture_nologin_screened.json")
		str, _ = nuxt.Create(string(f))
		n      = Nuxt{str}
	)

	{
		// Anonymous session has no user
		u, ok := n.User()
		assert.Equal(User{}, u)
		assert.False(ok)
	}

	// kamisama-day
	k := n.Radios()[34]

	eqs := []struct {
		in  interface{}
		out interface{}
	}{
		// radio
		{len(n.Radios()), 128},
		{k.Id(), 139},
		{k.Name(), "kamisama-day"},
		{k.Title(), "神様になったラジオ"},
		{k.HasBeenUpdated(), false},
		{len(k.Hosts()), 2},
		{k.Hosts()[0].Id(), 55},
		{k.Hosts()[1].Name(), "花江夏樹"},
		{len(k.Episodes()), 8},
		{k.Episodes()[0].Title(), "第12回"},
		{k.Episodes()[1].Title(), "第12回 おまけ"},
		// episode 7
		{k.Episodes()[6].Id(), 3114},
		{k.Episodes()[6].RadioId(), 139},
		{k.Episodes()[6].Title(), "第9回"},
		{k.Episodes()[6].Title(), "第9回"},
		{
			k.Episodes()[6].Poster(),
			"https://d3bzklg4lms4gh.cloudfront.net/program_info/image/default/production" +
				"/66/99/05f3c9402ca36cc3156dd50b7ab9aad298dd/image?v=1602579721",
		},
		{k.Episodes()[6].Guests(), []string{"重松千晴"}},
		{k.Episodes()[6].IsBonus(), false},
		{k.Episodes()[6].IsSticky(), false},
		{k.Episodes()[6].IsLatest(), false},
		{k.Episodes()[6].RequiresPremium(), true},
		{k.Episodes()[6].HasVideoStream(), false},
	}
	for _, eq := range eqs {
		assert.Equal(eq.out, eq.in)
	}

	{
		m, ok := k.Episodes()[6].Manifest()
		assert.False(ok)
		assert.Equal("", m, "Anonymous user is unable to access paid contents")
	}
	{
		// Radio time
		at, ok := k.JstUpdatedAt()
		name, offset := at.Zone()

		assert.True(ok)
		// year is variant
		assert.Equal(time.Month(3), at.Month())
		assert.Equal(19, at.Day())
		assert.Equal("UTC+9", name)
		assert.Equal(9*60*60, offset)
	}
	{

		// Episode time
		at, ok := k.Episodes()[6].JstUpdatedAt()
		name, offset := at.Zone()
		assert.True(ok)
		assert.Equal(time.Month(2), at.Month())
		assert.Equal(5, at.Day())
		assert.Equal("UTC+9", name)
		assert.Equal(9*60*60, offset)
	}
	{
		// 100man: no radio time, but contains episodes, use episode's time
		r := n.Radios()[0]
		assert.Nil(r.Raw.Updated)

		a, ok := r.JstUpdatedAt()
		assert.True(ok)
		b, _ := r.Episodes()[0].JstUpdatedAt()
		assert.Equal(b, a)
	}
	{
		// sorasara: no radio time, no episodes
		r := n.Radios()[3]
		assert.Nil(r.Raw.Updated)

		a, ok := r.JstUpdatedAt()
		assert.False(ok)
		assert.Equal(time.Time{}, a)
	}
	{
		// fujita: has videos
		f := n.Radios()[17]
		assert.True(f.Episodes()[0].HasVideoStream())
	}
}

func TestNuxtWithPremiumUser(t *testing.T) {
	var (
		assert = assert.New(t)
		f, _   = os.ReadFile("testdata/fixture_paid_screened.json")
		str, _ = nuxt.Create(string(f))
		n      = Nuxt{str}
	)

	u, ok := n.User()

	eqs := []struct {
		in  interface{}
		out interface{}
	}{
		{ok, true},
		{u.Email(), "hello@world"},
		{u.Id(), "0"},
		{
			u.FollowingPeople(),
			[]int{
				1377, 1044, 946, 889, 726, 645, 641, 590, 559, 429, 421, 284, 211, 136,
				114, 113, 105, 77, 66, 55, 396, 29,
			},
		},
		{
			u.FollowingRadios(),
			[]int{
				4, 10, 16, 17, 18, 29, 47, 54, 56, 65, 76, 77, 88, 89, 93, 118, 131, 136,
				139, 149, 156, 159,
			},
		},
		{u.PlaylistEpisodes(), []int{3676, 3677}},
	}
	for _, eq := range eqs {
		assert.Equal(eq.out, eq.in)
	}

	{
		// tsudaken: premium user can watch premium content
		assert.Equal("tsudaken", n.Radios()[17].Name())

		m, ok := n.Radios()[17].Episodes()[0].Manifest()
		assert.True(ok)
		assert.Equal("HAS_BEEN_SCREENED", m)
	}
}
