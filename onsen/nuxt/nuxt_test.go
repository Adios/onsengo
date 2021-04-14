package nuxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	{
		n, err := Create("")
		assert.Nil(n)
		assert.EqualError(err, "EOF", "Empty string causes error")
	}
	{
		n, err := Create("{\"ok\":\"}")
		assert.Nil(n)
		assert.EqualError(err, "unexpected EOF", "Invalid JSON causes error")
	}
	{
		n, err := Create("{}")
		assert.NoError(err)
		assert.Equal([]Program(nil), n.State.Programs.Programs.All, "Empty JSON creates entire valid structs")
	}
}

func TestCreateWithAnonymousUser(t *testing.T) {
	var (
		assert = assert.New(t)
		f, _   = os.Open("../testdata/fixture_nologin_screened.json")
	)
	defer f.Close()

	n, err := CreateFromReader(f)
	assert.NoError(err)
	assert.Nil(n.State.Signin, "Anonymous user doesn't have signin value")

	kami := n.State.Programs.Programs.All[34]

	equals := []struct {
		in       interface{}
		expected interface{}
	}{
		{len(n.State.Programs.Programs.All), 128},
		{kami.Id, 139},
		{kami.DirectoryName, "kamisama-day"},
		{kami.Title, "神様になったラジオ"},
		{kami.New, false},
		{*kami.Updated, "3/19"},
		{len(kami.Contents), 8},
		{
			kami.Performers,
			[]Performer{
				{55, "佐倉綾音"},
				{140, "花江夏樹"},
			},
		},
		{
			kami.Contents[1],
			Content{
				3677, "第12回 おまけ",
				true, false, true, "sound", true,
				139, "3/19", false,
				"https://d3bzklg4lms4gh.cloudfront.net/program_info/image/default/production" +
					"/66/99/05f3c9402ca36cc3156dd50b7ab9aad298dd/image?v=1602579721",
				nil, []string{},
			},
		},
	}

	for _, eq := range equals {
		assert.Equal(eq.expected, eq.in)
	}

	{
		// There are programs with nil updated value:
		// 100man
		assert.Nil(n.State.Programs.Programs.All[0].Updated)
		// sorasara
		assert.Nil(n.State.Programs.Programs.All[3].Updated)
	}
}

func TestCreateWithLogined(t *testing.T) {
	var (
		assert = assert.New(t)
		f, _   = os.Open("../testdata/fixture_paid_screened.json")
	)
	defer f.Close()

	n, err := CreateFromReader(f)
	assert.NoError(err)
	assert.NotNil(n.State.Signin, "Logined user have signin value")

	kami := n.State.Programs.Programs.All[30]

	equals := []struct {
		in       interface{}
		expected interface{}
	}{
		{len(n.State.Programs.Programs.All), 128},
		{kami.Id, 139},
		{kami.DirectoryName, "kamisama-day"},
		{kami.Title, "神様になったラジオ"},
		{kami.New, false},
		{*kami.Updated, "3/19"},
		{len(kami.Contents), 8},
		// Preimum user can access this content
		{*kami.Contents[1].StreamingUrl, "HAS_BEEN_SCREENED"},
	}

	for _, eq := range equals {
		assert.Equal(eq.expected, eq.in)
	}
}

func Example() {
	f, err := os.Open("../testdata/fixture_nologin_screened.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := CreateFromReader(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(n.State.Programs.Programs.All[34].Title)
	// Output: 神様になったラジオ
}
