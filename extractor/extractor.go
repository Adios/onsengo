// Package extractor is to deobfuscate the NUXT variable exported by nuxt.js.
package extractor

import (
	"errors"
	"regexp"

	"github.com/dop251/goja"
)

type Extractor interface {
	// Run the underlying JS vm to deobfuscate NUXT data, returns a JSON.stringify string if found
	Extract() (string, error)
}

type extractor struct {
	toextract string
	extracted *string
	vm        *goja.Runtime
}

// Take a string where you want to extract for NUXT data, usually an entire html content.
func New(str string) Extractor {
	return &extractor{
		toextract: str,
	}
}

func (e *extractor) Extract() (string, error) {
	if e.extracted == nil {
		data, err := e.getFirstObfuscatedData()
		if err != nil {
			return "", err
		}
		data, err = e.deobfuscate(data)
		if err != nil {
			return "", err
		}
		e.extracted = &data
	}
	return *e.extracted, nil
}

func (e *extractor) getVm() *goja.Runtime {
	if e.vm == nil {
		e.vm = goja.New()
	}
	return e.vm
}

func (e *extractor) getFirstObfuscatedData() (string, error) {
	re := regexp.MustCompile("<script>window.(__NUXT__=[^<]+;)</script>")
	m := re.FindStringSubmatch(e.toextract)

	if m == nil {
		return "", errors.New("NUXT data not found")
	}
	return m[1], nil
}

func (e *extractor) deobfuscate(str string) (string, error) {
	_, err := e.getVm().RunString(str)
	if err != nil {
		return "", err
	}
	res, err := e.getVm().RunString("JSON.stringify(__NUXT__)")
	if err != nil {
		return "", err
	}
	return res.Export().(string), nil
}
