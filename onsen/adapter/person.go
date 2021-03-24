package adapter

import (
	"github.com/adios/onsengo/onsen/nuxt"
)

// Person wraps an instance of nuxt.Performer to transform output on the fly.
type Person interface {
	PersonId() uint
	Name() string
}

type person struct {
	raw *nuxt.Performer
}

func (p person) PersonId() uint {
	return uint(p.raw.Id)
}

func (p person) Name() string {
	return p.raw.Name
}

func NewPerson(p *nuxt.Performer) Person {
	if p == nil {
		panic("Cannot be nil")
	}
	return person{p}
}
