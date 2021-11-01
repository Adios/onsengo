package onsen

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/adios/onsengo/onsen/nuxt"
)

func TestMain(m *testing.M) {
	SetRefDate("2021-10-29")
	os.Exit(m.Run())
}

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

func TestGuessJstTimeWithNow(t *testing.T) {
	mem := guessRefTime
	defer func() { guessRefTime = mem }()

	{
		SetRefDate("2021-03-24")
		tm, ok := GuessJstTimeWithNow("3/25")
		assert.True(t, ok)
		assert.Equal(t, "2020-03-25 00:00:00 +0900 UTC+9", tm.String())
	}
	{
		SetRefDate("2020-03-26")
		tm, ok := GuessJstTimeWithNow("3/25")
		assert.True(t, ok)
		assert.Equal(t, "2020-03-25 00:00:00 +0900 UTC+9", tm.String())
	}
	{
		SetRefDate("2019-03-24")
		tm, ok := GuessJstTimeWithNow("3/25")
		assert.True(t, ok)
		assert.Equal(t, "2018-03-25 00:00:00 +0900 UTC+9", tm.String())
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

	// radionyan
	k := n.Radios()[7]

	eqs := []struct {
		in  interface{}
		out interface{}
	}{
		// radio
		{len(n.Radios()), 141},
		{k.Id(), 202},
		{k.Name(), "radionyan"},
		{k.Title(), "月とライカと吸血姫 ～アーニャ・シモニャン・ラジオニャン！～"},
		{k.HasBeenUpdated(), false},
		{len(k.Hosts()), 1},
		{k.Hosts()[0].Id(), 1189},
		{k.Hosts()[0].Name(), "木野日菜"},
		{len(k.Episodes()), 6},
		{k.Episodes()[0].Title(), "第3回"},
		{k.Episodes()[1].Title(), "第3回 おまけ"},
		// episode 5
		{k.Episodes()[4].Id(), 6006},
		{k.Episodes()[4].RadioId(), 202},
		{k.Episodes()[4].Title(), "第1回"},
		{
			k.Episodes()[4].Poster(),
			"https://d3bzklg4lms4gh.cloudfront.net/program_info/image/default/production" +
				"/5b/6e/2a28979284885466f12fcc07f0a311736e29/image?v=1633683939",
		},
		{k.Episodes()[4].Guests()[0].Id(), 574},
		{k.Episodes()[4].Guests()[0].Name(), "内山昂輝"},
		{k.Episodes()[4].IsBonus(), false},
		{k.Episodes()[4].IsSticky(), false},
		{k.Episodes()[4].IsLatest(), false},
		{k.Episodes()[4].RequiresPremium(), true},
		{k.Episodes()[4].HasVideoStream(), false},
	}
	for _, eq := range eqs {
		assert.Equal(eq.out, eq.in)
	}

	{
		m, ok := k.Episodes()[4].Manifest()
		assert.False(ok)
		assert.Equal("", m, "Anonymous user is unable to access paid contents")
	}
	{
		// Radio time
		at, ok := k.JstUpdatedAt()
		assert.True(ok)
		assert.Equal("2021-10-22 00:00:00 +0900 UTC+9", at.String())
	}
	{
		// Episode time
		at, ok := k.Episodes()[4].JstUpdatedAt()
		assert.True(ok)
		assert.Equal("2021-09-17 00:00:00 +0900 UTC+9", at.String())
	}
	{
		// fujita: has videos
		f := n.Radios()[9]
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
		{u.Uid(), "0"},
		{
			u.FollowingPeople(),
			[]int{
				1544, 1044, 946, 889, 726, 645, 641, 590, 559, 429, 421, 284, 211, 136,
				114, 113, 105, 77, 66, 55, 396, 30, 1699, 493, 29,
			},
		},
		{
			u.FollowingRadios(),
			[]int{
				4, 10, 16, 17, 18, 29, 46, 47, 54, 65, 76, 77, 88, 89, 93, 96, 118, 134,
				179, 185, 193, 201, 202, 203, 216, 221,
			},
		},
		{
			u.PlaylistEpisodes(),
			[]int{
				6525, 6527, 6528, 6582, 6600, 6601, 6594, 6609, 6610,
			},
		},
	}
	for _, eq := range eqs {
		assert.Equal(eq.out, eq.in)
	}

	{
		// tsudaken: premium user can watch premium content
		assert.Equal("fujita", n.Radios()[9].Name())

		m, ok := n.Radios()[9].Episodes()[1].Manifest()
		assert.True(ok)
		assert.Equal("HAS_BEEN_SCREENED", m)
	}
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	{
		o, err := Create("")
		assert.Error(err)
		assert.Nil(o)
	}
	{
		o, err := Create("...<script>window.__NUXT__=one</script><script>window.__NUXT__=two;</script>...")
		assert.Error(err)
		assert.Nil(o)
	}
}

func TestOnsenRadio(t *testing.T) {
	var (
		assert = assert.New(t)
		f, _   = os.ReadFile("testdata/fixture_nologin_screened.html")
	)

	{
		o, _ := Create(string(f))
		r, ok := o.Radio(nil)
		assert.False(ok)
		assert.Equal(r, Radio{})
	}
	{
		o, _ := Create(string(f))
		assert.Nil(o.cache.r)
		r, ok := o.Radio(202)
		assert.True(ok)
		assert.Equal(o.Radios()[7], r)
		assert.NotNil(o.cache.r)
	}
	{
		o, _ := Create(string(f))
		assert.Nil(o.cache.r)
		a, ok := o.Radio("radionyan")
		assert.True(ok)
		assert.Equal(o.Radios()[7], a)
		assert.NotNil(o.cache.r)
		b, _ := o.Radio("radionyan")
		assert.Equal(b, a)
	}
}

func TestOnsenEpisode(t *testing.T) {
	var (
		assert = assert.New(t)
		f, _   = os.ReadFile("testdata/fixture_nologin_screened.html")
	)

	{
		o, _ := Create(string(f))
		assert.Nil(o.cache.e)
		r, ok := o.Episode(-1)
		assert.False(ok)
		assert.Equal(r, Episode{})
		assert.NotNil(o.cache.e)
	}
	{
		o, _ := Create(string(f))
		assert.Nil(o.cache.e)
		e, ok := o.Episode(6505)
		assert.True(ok)
		assert.Equal(o.Radios()[7].Episodes()[0], e)
		assert.NotNil(o.cache.e)
	}
}
