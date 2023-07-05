package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)





func blockIndefinitely(w http.ResponseWriter, r*http.Request) {
	select{}
}

func TestBlockIndefinitely(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))
	_, _ = http.Get(ts.URL)
	t.Fatal("Client did not indefinitely block")
}

