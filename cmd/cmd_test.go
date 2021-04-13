package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var (
		assert   = assert.New(t)
		server   = httptest.NewServer(server(t))
		out, err = new(bytes.Buffer), new(bytes.Buffer)
	)
	defer server.Close()

	root.cmd.SetOut(out)
	root.cmd.SetErr(err)
	root.out, root.err = out, err

	root.cmd.SetArgs([]string{"ls", "--backend", "file:///"})
	assert.EqualError(Execute(), "Create: NUXT pattern not matched")

	root.cmd.SetArgs([]string{"ls", "--backend", server.URL})
	out.Reset()
	err.Reset()
	{
		f, _ := os.ReadFile("testdata/expected_ls.txt")
		assert.NoError(Execute())
		assert.Equal(string(f), out.String())
	}

	root.cmd.SetArgs([]string{"ls", "-r", "--backend", server.URL})
	out.Reset()
	{
		f, _ := os.ReadFile("testdata/expected_ls_recursive.txt")
		assert.NoError(Execute())
		assert.Equal(string(f), out.String())
	}

	root.cmd.SetArgs([]string{"ls", "fujita", "gurepap", "nosuchradio", "--backend", server.URL})
	out.Reset()
	{
		f, _ := os.ReadFile("testdata/expected_ls_single.txt")
		assert.NoError(Execute())
		assert.Equal(string(f), out.String())
		assert.Equal("nosuchradio: not found\n", err.String())
	}
}

func server(t *testing.T) http.Handler {
	f, _ := os.ReadFile("testdata/fixture_nologin_screened.html")

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, ua, req.Header.Get("User-Agent"))
		w.Write(f)
	})
}
