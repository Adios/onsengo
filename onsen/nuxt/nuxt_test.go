package nuxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	nuxt, err := Parse("")

	assert.Nil(t, nuxt)
	assert.NotNil(t, err)
}

func TestInvalidNuxtObject(t *testing.T) {
	nuxt, err := Parse("{}")

	assert.Nil(t, err)

	all := nuxt.State.Programs.Programs.All

	assert.Equal(t, len(all), 0)
}

func TestNonlogined(t *testing.T) {
	content, err := os.ReadFile("fixture_nologin_screened.json")
	if err != nil {
		panic(err)
	}

	nuxt, err := Parse(string(content))

	assert.Nil(t, err)
	assert.Nil(t, nuxt.Error)
	assert.Equal(t, nuxt.RoutePath, "/")
	assert.Nil(t, nuxt.State.Signin)

	all := nuxt.State.Programs.Programs.All

	assert.Equal(t, len(all), 128)

	kami := all[34]

	assert.Equal(t, kami.Id, ProgramId(139))
	assert.Equal(t, kami.DirectoryName, "kamisama-day")
	assert.Equal(t, kami.Title, "神様になったラジオ")
	assert.Equal(t, kami.New, false)
	assert.Equal(t, kami.Updated, "3/19")
	assert.Equal(t, kami.Performers, []Performer{
		{55, "佐倉綾音"},
		{140, "花江夏樹"},
	})

	assert.Equal(t, len(kami.Contents), 8)
	assert.Equal(t, kami.Contents[1], Content{
		ContentId(3677), "第12回 おまけ",
		true, false, true, "sound", true,
		ProgramId(139), "3/19", false,
		"https://d3bzklg4lms4gh.cloudfront.net/program_info/image/default/production" +
			"/66/99/05f3c9402ca36cc3156dd50b7ab9aad298dd/image?v=1602579721",
		nil, []string{},
	})
}

func TestPaidMember(t *testing.T) {
	content, err := os.ReadFile("fixture_paid_screened.json")
	if err != nil {
		panic(err)
	}

	nuxt, err := Parse(string(content))

	assert.Nil(t, err)
	assert.Nil(t, nuxt.Error)
	assert.Equal(t, nuxt.RoutePath, "/")
	assert.NotNil(t, nuxt.State.Signin)

	all := nuxt.State.Programs.Programs.All

	assert.Equal(t, len(all), 128)

	kami := all[30]

	assert.Equal(t, kami.Id, ProgramId(139))
	assert.Equal(t, kami.DirectoryName, "kamisama-day")
	assert.Equal(t, kami.Title, "神様になったラジオ")
	assert.Equal(t, kami.New, false)
	assert.Equal(t, kami.Updated, "3/19")
	assert.Equal(t, kami.Performers, []Performer{
		{55, "佐倉綾音"},
		{140, "花江夏樹"},
	})

	assert.Equal(t, len(kami.Contents), 8)
	assert.Equal(t, *kami.Contents[1].StreamingUrl, "HAS_BEEN_SCREENED")
}

func Example() {
	content, err := os.ReadFile("fixture_nologin_screened.json")
	if err != nil {
		panic(err)
	}

	nuxt, err := Parse(string(content))
	if err != nil {
		panic(err)
	}

	fmt.Println(nuxt.State.Programs.Programs.All[34].Title)
	// Output: 神様になったラジオ
}
