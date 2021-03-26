package onsen

import (
	"errors"

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
