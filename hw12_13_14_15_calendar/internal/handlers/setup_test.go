package handlers

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"

	"github.com/go-chi/chi/v5" //nolint:typecheck
)

var (
	router http.Handler
	store  *memorystorage.Storage
)

func TestMain(m *testing.M) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal("failed to load location")
	}
	time.Local = loc

	logg, err := logger.New("debug", true)
	if err != nil {
		log.Fatal("failed to make a logger")
	}

	store = memorystorage.New()
	calendar := app.New(logg, store)
	handlers := NewHandlers(calendar, logg)

	mux := chi.NewRouter() //nolint:typecheck
	mux.Post("/events", handlers.CreateEvent)
	mux.Put("/events/{eventID}", handlers.UpdateEvent)
	mux.Delete("/events/{eventID}", handlers.DeleteEvent)
	mux.Get("/events", handlers.ListEvent)
	mux.Get("/events/{eventID}", handlers.GetEvent)
	router = mux

	os.Exit(m.Run())
}
