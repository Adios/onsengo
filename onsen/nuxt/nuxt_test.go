package nuxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyString(t *testing.T) {
	nuxt, err := From("")

	assert.Nil(t, nuxt)
	assert.NotNil(t, err)
}

func TestInvalidJsonString(t *testing.T) {
	nuxt, err := From("{}")

	assert.Nil(t, err)

	all := nuxt.State.Programs.Programs.All

	assert.Equal(t, 0, len(all))
}

func TestUsingAnonymous(t *testing.T) {
	f, err := os.Open("testdata/fixture_nologin_screened.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	nuxt, err := FromReader(f)

	assert.Nil(t, err)
	assert.Nil(t, nuxt.Error)
	assert.Equal(t, "/", nuxt.RoutePath)
	assert.Nil(t, nuxt.State.Signin)

	all := nuxt.State.Programs.Programs.All

	assert.Equal(t, 128, len(all))

	kami := all[34]

	assert.Equal(t, uint(139), kami.Id)
	assert.Equal(t, "kamisama-day", kami.DirectoryName)
	assert.Equal(t, "神様になったラジオ", kami.Title)
	assert.Equal(t, false, kami.New)
	assert.Equal(t, "3/19", *kami.Updated)
	assert.Equal(
		t,
		[]Performer{
			{55, "佐倉綾音"},
			{140, "花江夏樹"},
		},
		kami.Performers,
	)

	assert.Equal(t, 8, len(kami.Contents))
	assert.Equal(
		t,
		Content{
			3677, "第12回 おまけ",
			true, false, true, "sound", true,
			139, "3/19", false,
			"https://d3bzklg4lms4gh.cloudfront.net/program_info/image/default/production" +
				"/66/99/05f3c9402ca36cc3156dd50b7ab9aad298dd/image?v=1602579721",
			nil, []string{},
		},
		kami.Contents[1],
	)
}

func TestUsingPaidMember(t *testing.T) {
	f, err := os.Open("testdata/fixture_paid_screened.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	nuxt, err := FromReader(f)

	assert.Nil(t, err)
	assert.Nil(t, nuxt.Error)
	assert.Equal(t, "/", nuxt.RoutePath)
	assert.NotNil(t, nuxt.State.Signin)

	all := nuxt.State.Programs.Programs.All

	assert.Equal(t, 128, len(all))

	kami := all[30]

	assert.Equal(t, uint(139), kami.Id)
	assert.Equal(t, "kamisama-day", kami.DirectoryName)
	assert.Equal(t, "神様になったラジオ", kami.Title)
	assert.Equal(t, false, kami.New)
	assert.Equal(t, "3/19", *kami.Updated)
	assert.Equal(
		t,
		[]Performer{
			{55, "佐倉綾音"},
			{140, "花江夏樹"},
		},
		kami.Performers,
	)

	assert.Equal(t, 8, len(kami.Contents))
	assert.Equal(t, "HAS_BEEN_SCREENED", *kami.Contents[1].StreamingUrl)
}

func TestThereAreNilUpdatedTimes(t *testing.T) {
	f, err := os.Open("testdata/fixture_nologin_screened.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	nuxt, err := FromReader(f)

	assert.Nil(t, err)

	all := nuxt.State.Programs.Programs.All

	p1 := all[0]
	p2 := all[3]

	assert.Equal(t, "100man", p1.DirectoryName)
	assert.Nil(t, p1.Updated)

	assert.Equal(t, "sorasara", p2.DirectoryName)
	assert.Nil(t, p2.Updated)
}

func Example() {
	f, err := os.Open("testdata/fixture_nologin_screened.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	nuxt, err := FromReader(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(nuxt.State.Programs.Programs.All[34].Title)
	// Output: 神様になったラジオ
}
