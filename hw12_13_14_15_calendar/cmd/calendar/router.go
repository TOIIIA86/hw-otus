package main

import (
	"errors"
	"os"

	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/handlers"
	internalrouter "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/router"
)

const routerChi = "chi"

var undefinedRouterType = errors.New("undefined router type")

func NewRouter(logfile *os.File, handlers *handlers.Handlers, config Config) (app.Router, error) {
	var router app.Router

	switch config.App.RouterType {
	case routerChi:
		router = internalrouter.NewChiRouter(logfile, handlers)
	default:
		return nil, undefinedRouterType
	}

	return router, nil
}
