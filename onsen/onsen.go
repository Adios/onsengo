package onsen

import (
	"fmt"
	"strconv"

	"github.com/adios/onsengo/js/expression"
	"github.com/adios/onsengo/onsen/nuxt"
	"github.com/adios/onsengo/onsen/nuxt/decorator"
)

type Onsen struct {
	decorator.Decorator

	NuxtJson string
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

func Create(htmlstr string) (*Onsen, error) {
	js, ok := FindNuxtExpression(htmlstr)
	if !ok {
		return nil, fmt.Errorf("NUXT pattern not matched")
	}

	jsonstr, err := expression.From(js).Stringify()
	if err != nil {
		return nil, err
	}

	n, err := nuxt.From(jsonstr)
	if err != nil {
		return nil, err
	}

	return &Onsen{
		Decorator: decorator.Decorator{
			Raw: n,
		},
		NuxtJson: jsonstr,
	}, nil
}
