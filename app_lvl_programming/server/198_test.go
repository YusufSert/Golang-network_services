package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"
)

type Methods map[string]http.Handler

func (h Methods) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(ioutil.Discard, r)
		_ = r.Close()
	}(r.Body)

	if handler, ok := h[r.Method]; ok {
		if handler == nil {
			http.Error(w, "Internal server error",
			http.StatusInternalServerError)
		} else  {
			handler.ServeHTTP(w, r)
		}

		return
	}


	w.Header().Add("Allow", h.allowedMethods())
	if r.Method != http.MethodOptions {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func (h Methods) allowedMethods() string {
	a := make([]string, 0, len(h))

	for k := range h {
		a = append(a, k)
	}
	
	sort.Strings(a)
	return strings.Join(a, ", ")
}


func DefaultMethodsHandler() http.Handler {
	return Methods{
		http.MethodGet: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("Hello, friend!"))
			},
		),

		http.MethodPost: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Internal server error",
					http.StatusInternalServerError)
					return
				}

				_, _ = fmt.Fprintf(w, "Hello, %s!",
				html.EscapeString(string(b)))
			},
		),
	}
}

func TestSimpleHTTPServer(t *testing.T) {
	srv := &http.Server{
		Addr: "127.0.0.1:8081",
		Handler: http.TimeoutHandler(
			DefaultMethodsHandler(), 2*time.Minute, ""),
		IdleTimeout: 5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}

	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	

	err = srv.Serve(l)
	if err != http.ErrServerClosed {
		t.Error(err)
	}
}

