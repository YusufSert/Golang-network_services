package main

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	//Create a listener on a random port.
	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		defer func() { done <- struct{}{} }()

		for {
			conn, err := listener.Accept() // This will block until detect a incoming connection and completes handshake between server and client
			if err != nil {
				t.Log(err)
				return
			}

			// Handle conn in a new goroutine
			go func(c net.Conn) {
				defer func() {
					c.Close() // sends fin packet to client
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn) // End of goroutine
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String()) // sends ack to server
	if err != nil {
		t.Fatal(err)
	}

	conn.Close() // sends fin packet to server
	<-done
	listener.Close() // unblocks the listener accept method
	<-done
}
