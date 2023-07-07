package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)



func main() {
	
	srv := &http.Server{
		Addr: "127.0.0.1:8081",
		Handler: Middleware(http.HandlerFunc(GetHandler)),
	}
	
	l, err := net.Listen("tcp", srv.Addr)

	if err != nil {
		log.Fatal(err)
	}

	err = srv.Serve(l)
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}


}



func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			log.Println(r.Method)

			next.ServeHTTP(w, r)
		},
	)

}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(ioutil.Discard, r)
		_ = r.Close()

	}(r.Body)

	w.WriteHeader(http.StatusPaymentRequired)
	w.Write([]byte("Middleware"))

}


