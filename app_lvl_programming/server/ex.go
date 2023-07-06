package main

import (
	"database/sql"
	"log"
	"net/http"
)




type Handlers struct {
	db *sql.DB
	log *log.Logger
}

func dbHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request ) {
			_ = db.Ping()
		},
	)
}
func (h *Handlers) Handler1() http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			err := h.db.Ping()
			if err != nil {
				h.log.Printf("db ping: %v", err)
			}
			// do something with the databse here
		},

	)

}

