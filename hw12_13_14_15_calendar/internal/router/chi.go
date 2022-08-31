package router

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5" //nolint:typecheck
	"github.com/go-chi/chi/v5/middleware"
)

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Application interface {
	GetName() string
}

type Handlers interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
	ListEvent(w http.ResponseWriter, r *http.Request)
	GetEvent(w http.ResponseWriter, r *http.Request)
}

func NewChiRouter(logfile *os.File, handlers Handlers) Router {
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(logfile, "", log.LstdFlags), NoColor: true})

	mux := chi.NewRouter() //nolint:typecheck
	mux.Use(middleware.Logger)

	mux.Post("/events", handlers.CreateEvent)
	mux.Put("/events/{eventID}", handlers.UpdateEvent)
	mux.Delete("/events/{eventID}", handlers.DeleteEvent)
	mux.Get("/events", handlers.ListEvent)
	mux.Get("/events/{eventID}", handlers.GetEvent)
	return mux
}
