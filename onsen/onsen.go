// Package onsen embeds a Nuxt decorator acting as a front-end to clients
package onsen

import (
	"fmt"
	"strconv"

	"github.com/adios/onsengo/js/expression"
	"github.com/adios/onsengo/onsen/nuxt"
	"github.com/adios/onsengo/onsen/nuxt/decorator"
)

type RadioShow = decorator.RadioShow

type Onsen struct {
	decorator.Decorator

	NuxtJson   string
	radioCache *RadioCache
}

// Accepts either a RadioShowId or RadioShowName, returns a *RadioShow if found.
// It creates a cache of all the radio shows when being invoked for first time.
func (o *Onsen) RadioShow(id Identifier) (r *RadioShow, ok bool) {
	switch v := id.(type) {
	case RadioShowName:
		return o.RadioCache().NameShow(v.Identify())
	case RadioShowId:
		return o.RadioCache().IdShow(v.Identify())
	default:
		return nil, false
	}
}

func (o *Onsen) RadioCache() *RadioCache {
	if o.radioCache == nil {
		o.radioCache = createRadioCache(o)
	}
	return o.radioCache
}

type Identifier interface {
	Identify() string
}

type RadioShowId uint

func (id RadioShowId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func (id RadioShowId) Identify() string {
	return id.String()
}

type RadioShowName string

func (n RadioShowName) Identify() string {
	return string(n)
}

type EpisodeId uint

func (id EpisodeId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func (id EpisodeId) Identify() string {
	return id.String()
}

type PersonId uint

func (id PersonId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func (id PersonId) Identify() string {
	return id.String()
}

type UserId string

func (id UserId) Identify() string {
	return string(id)
}

type RadioCache struct {
	idShow   map[string]*RadioShow
	nameShow map[string]*RadioShow
}

func (c *RadioCache) IdShow(id string) (r *RadioShow, ok bool) {
	r, ok = c.idShow[id]
	if !ok {
		return nil, false
	}
	return r, true
}

func (c *RadioCache) NameShow(name string) (r *RadioShow, ok bool) {
	r, ok = c.nameShow[name]
	if !ok {
		return nil, false
	}
	return r, true
}

type radioShowWalker interface {
	EachRadioShow(func(RadioShow))
}

func createRadioCache(w radioShowWalker) *RadioCache {
	cache := RadioCache{
		idShow:   map[string]*RadioShow{},
		nameShow: map[string]*RadioShow{},
	}

	w.EachRadioShow(func(r RadioShow) {
		cache.idShow[RadioShowId(r.RadioShowId()).String()] = &r
		cache.nameShow[r.Name()] = &r
	})

	return &cache
}

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
