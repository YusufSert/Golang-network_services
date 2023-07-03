package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanOut(t *testing.T) {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(10*time.Second))

	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()

	go func() {
		// Only accepting a single connection
		conn, err := listener.Accept()

		if err == nil {
			conn.Close()
		}

	}()

	dial := func(ctx context.Context, address string, response chan int,
		id int, wg *sync.WaitGroup) {

		defer wg.Done()
		var d net.Dialer

		c, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			nErr, ok := err.(net.Error)
			fmt.Println(nErr.Error(), ok)
			return
		}
		fmt.Println(c.RemoteAddr(), id, c.LocalAddr())
		c.Close()

		select {
		case <-ctx.Done():
			fmt.Println("done", id)
		case response <- id:
			fmt.Println("res", id)
		}
	}

	res := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	response := <-res
	cancel()
	wg.Wait()
	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %s", ctx.Err())
	}
	t.Logf("dialer %d retriwed the resource", response)
}
