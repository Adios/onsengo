package adapter

import (
	"github.com/adios/onsengo/onsen/adapter"
	"github.com/adios/onsengo/onsen/nuxt"
)

type person struct {
	raw *nuxt.Performer
}

func (p person) PersonId() uint {
	return uint(p.raw.Id)
}

func (p person) Name() string {
	return p.raw.Name
}

// Person wraps an instance of nuxt.Performer to transform output on the fly.
func NewPerson(p *nuxt.Performer) adapter.Person {
	if p == nil {
		panic("Cannot be nil")
	}
	return person{p}
}
