package onsen

import (
	"regexp"
)

// Returns a byte slice to the capture of first appeared NUXT pattern:
//   <script>window.__NUXT__=([^<]+);</script>
func FindNuxtExpression(html []byte) (expr []byte, ok bool) {
	re := regexp.MustCompile("<script>window.__NUXT__=([^<]+);</script>")

	m := re.FindSubmatch(html)
	if m == nil {
		return nil, false
	}
	return m[1], true
}
