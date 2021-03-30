package onsen

import (
	"errors"
	"strconv"

	"github.com/adios/onsengo/js/expression"
	"github.com/adios/onsengo/onsen/adapter"
	"github.com/adios/onsengo/onsen/html"
	"github.com/adios/onsengo/onsen/nuxt"
	nuxtadapter "github.com/adios/onsengo/onsen/nuxt/adapter"
)

type Onsen interface {
	RadioShows() []adapter.RadioShow
}

type onsen struct {
	data adapter.Adapter
}

func (o *onsen) RadioShows() []adapter.RadioShow {
	return o.data.RadioShows()
}

type RadioShowId uint

func (id RadioShowId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type EpisodeId uint

func (id EpisodeId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type PersonId uint

func (id PersonId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type UserId string

func Create(htmlstr string) (Onsen, error) {
	h := html.FindNuxtExpression(htmlstr)
	if h == nil {
		return nil, errors.New("No known js nuxt object found in the html")
	}

	jsonstr, err := expression.New(*h).Stringify()
	if err != nil {
		return nil, err
	}

	nu, err := nuxt.Parse(jsonstr)
	if err != nil {
		return nil, err
	}

	a := nuxtadapter.NewAdapter(nu)

	return &onsen{
		data: a,
	}, nil
}
