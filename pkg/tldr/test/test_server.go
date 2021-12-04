package tldrtest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

const tldrZipFilename = "tldr.zip"

func tmpDir() string {
	return "/tmp/tldrtest"
}

type Informer interface {
	ServerURL() string
	TldrZipURL() string
	Close()
}

type Starter interface {
	Start() Informer
}

type TestServer struct {
	testServer *httptest.Server
}

// NewServer returns an instance that can start server.
// Alfter starting server, server information wll be available with methods
func NewServer() Starter {
	return &TestServer{
		testServer: newTldrRepositoryServer(),
	}
}

func (s *TestServer) Start() Informer {
	s.testServer.Start()
	return s
}

func (s *TestServer) Close() {
	s.testServer.Close()
}

func (s *TestServer) ServerURL() string {
	return s.testServer.URL
}

func (s *TestServer) TldrZipURL() string {
	return s.ServerURL() + "/" + tldrZipFilename
}

func newTldrRepositoryServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, tldrZipFilename) {
			fmt.Fprintf(w, "hello")
			return
		}

		zipPath := filepath.Join(tmpDir(), tldrZipFilename)
		if _, err := os.Stat(zipPath); err != nil {
			panic(err)
		}

		f, err := os.Open(zipPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			panic(err)
		}
	})

	return httptest.NewUnstartedServer(mux)
}
