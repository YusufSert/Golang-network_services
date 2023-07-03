package main

import (
	"net"
	"testing"
)

func TestListener2(t *testing.T) {
	listener, err := net.Listen("tpc", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	for {
		conn, _ := listener.Accept()

		go func(c net.Conn) {
			defer c.Close()
			// Your code would handle the connection here.
		}(conn)
	}
}
