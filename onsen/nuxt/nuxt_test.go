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

	chosen := n.State.Programs.Programs.All[7]

	equals := []struct {
		in       interface{}
		expected interface{}
	}{
		{len(n.State.Programs.Programs.All), 141},
		{chosen.Id, 202},
		{chosen.DirectoryName, "radionyan"},
		{chosen.Title, "月とライカと吸血姫 ～アーニャ・シモニャン・ラジオニャン！～"},
		{chosen.New, false},
		{*chosen.Updated, "10/22"},
		{len(chosen.Contents), 6},
		{
			chosen.Performers,
			[]Performer{
				{1189, "木野日菜"},
			},
		},
		{
			chosen.Contents[1],
			Content{
				6506, "第3回 おまけ",
				true, false, true, "sound", true,
				202, "10/22", false,
				"https://d3bzklg4lms4gh.cloudfront.net/program_info/image/default/production" +
					"/5b/6e/2a28979284885466f12fcc07f0a311736e29/image?v=1633683939",
				nil, []Performer{},
			},
		},
	}

	for _, eq := range equals {
		assert.Equal(eq.expected, eq.in)
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

	chosen := n.State.Programs.Programs.All[7]

	equals := []struct {
		in       interface{}
		expected interface{}
	}{
		{len(n.State.Programs.Programs.All), 141},
		{chosen.Id, 202},
		{chosen.DirectoryName, "radionyan"},
		{chosen.Title, "月とライカと吸血姫 ～アーニャ・シモニャン・ラジオニャン！～"},
		{chosen.New, false},
		{*chosen.Updated, "10/22"},
		{len(chosen.Contents), 6},
		// Preimum user can access this content
		{*chosen.Contents[1].StreamingUrl, "HAS_BEEN_SCREENED"},
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

	fmt.Println(n.State.Programs.Programs.All[7].Title)
	// Output: 月とライカと吸血姫 ～アーニャ・シモニャン・ラジオニャン！～
}
