package rest

import (
	"log"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"net/http"
)

type Server struct {}

func (s Server) Run() {
	log.Printf("[INFO] activate rest server")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/v1", func (r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}