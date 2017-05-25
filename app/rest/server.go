package rest

import (
	"log"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"net/http"
	"encoding/json"
	"github.com/asnelzin/pomodoro/app/models"
	"github.com/asnelzin/pomodoro/app/vk"
)

type Server struct {
	Bot                 *vk.Bot
	SecretKey           string
	ConfirmationMessage string
}

func (s Server) Run() {
	log.Printf("[INFO] activate rest server")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		r.Post("/callback", s.handleCallback)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

// POST /v1/callback
func (s Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	event := models.Event{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {

	case "confirmation":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s.ConfirmationMessage))

	case "message_new":
		if !CheckSecretMatch(event.SecretKey, s.SecretKey) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err = s.Bot.HandleMessage(&event.Object)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}