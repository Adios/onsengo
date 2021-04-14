// Package cmd implements the onsengo command.
package cmd

import (
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/adios/onsengo/onsen"
)

const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:87.0) Gecko/20100101 Onsengo/1.0"

var root = ctx{
	cmd: &cobra.Command{
		Use:   "onsengo",
		Short: "List and browse onsen.ag radio shows",
		Long: `
onsengo♨ is a program which allows browsing radio shows on https://onsen.ag.
`,
	},
}

func Execute() error {
	return root.cmd.Execute()
}

func init() {
	pf := root.cmd.PersistentFlags()

	pf.StringVar(&root.backend, "backend", "https://onsen.ag/", "set backend, file:// is supported")
	pf.StringVarP(&root.session, "session", "s", "", "set session")
}

type ctx struct {
	// To test or interpret an archived html, e.g.: file:///full/path/to/the/onsen/index.html
	backend string
	// You can find the id from "_session_id=SESSION_ID" in the browser's cookie.
	session string
	cmd     *cobra.Command

	// for onsen/pprint output
	out io.Writer
	err io.Writer

	// for onsen
	hc *http.Client
	oo *onsen.Onsen
}

func (c *ctx) client() *http.Client {
	if c.hc == nil {
		t := &http.Transport{}
		// for file://
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		c.hc = &http.Client{Transport: t}
	}
	return c.hc
}

func (c *ctx) onsen() (*onsen.Onsen, error) {
	if c.oo == nil {
		req, err := http.NewRequest("GET", c.backend, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("User-Agent", ua)
		if c.session != "" {
			req.Header.Add("Cookie", "_session_id="+c.session)
		}

		resp, err := c.client().Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		o, err := onsen.Create(string(b))
		if err != nil {
			return nil, err
		}
		c.oo = o
	}
	return c.oo, nil
}

func (c *ctx) outw() io.Writer {
	if c.out == nil {
		c.out = os.Stdout
	}
	return c.out
}

func (c *ctx) errw() io.Writer {
	if c.err == nil {
		c.err = os.Stderr
	}
	return c.err
}
