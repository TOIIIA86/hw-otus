package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Server       *http.Server
	Logger       Logger
	Router       Router
}

func NewServer(host, port string, readTimeout, writeTimeout, idleTimeout time.Duration, logger Logger, router Router) *Server {
	return &Server{
		Host:         host,
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Logger:       logger,
		Router:       router,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var ch = make(chan error)

	s.Logger.Info("Starting http server on " + s.getAddr())

	go func() {
		s.Server = &http.Server{
			Addr:         s.getAddr(),
			Handler:      s.Router,
			ReadTimeout:  s.ReadTimeout,
			WriteTimeout: s.WriteTimeout,
			IdleTimeout:  s.IdleTimeout,
		}
		ch <- s.Server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-ch:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *Server) getAddr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
