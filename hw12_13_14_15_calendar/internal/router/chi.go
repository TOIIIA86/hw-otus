package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Application interface {
	GetName() string
}

func NewChiRouter(logfile *os.File, app Application) Router {
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(logfile, "", log.LstdFlags), NoColor: true})

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("Hello, %s!", app.GetName())))
	})
	return mux
}
