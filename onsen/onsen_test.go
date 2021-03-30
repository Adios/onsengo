package onsen

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdConversion(t *testing.T) {
	assert.Equal(t, "0", PersonId(0).String())
	assert.Equal(t, "55", PersonId(55).String())
	assert.Equal(t, "555", RadioShowId(555).String())
	assert.Equal(t, "5555", EpisodeId(5555).String())
}

func TestCreateEmpty(t *testing.T) {
	o, err := Create("")

	assert.Equal(t, fmt.Errorf("NUXT pattern not matched"), err)
	assert.Nil(t, o)
}

func TestOK(t *testing.T) {
	content, err := os.ReadFile("testdata/fixture_nologin_screened.html")
	if err != nil {
		panic(err)
	}

	o, err := Create(string(content))

	assert.Nil(t, err)
	assert.Equal(t, 128, len(o.RadioShows()))
}
