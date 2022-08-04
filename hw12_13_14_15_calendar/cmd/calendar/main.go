package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	sqlstorage "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/storage/sql" //nolint:typecheck

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", ".env", "Path to configuration file")
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx := context.Background()

	// loading .env
	if err := godotenv.Load(configFile); err != nil {
		panic(err)
	}

	// populating config
	config := NewConfig()
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	// creating logger
	logg, err := logger.New(config.Logger.Level, config.Logger.Development)
	if err != nil {
		panic(err)
	}

	// creating storage
	storage, err := NewStorage(ctx, config)
	if err != nil {
		logg.Fatal("failed to create a storage: " + err.Error())
	}
	if s, ok := storage.(sqlstorage.Storage); ok {
		defer s.Close(ctx)
	}

	// creating application
	calendar := app.New(config.App.Name, logg, storage)

	// creating requests log file
	if err := os.Mkdir("./logs", 0755); err != nil && err == os.ErrExist {
		logg.Fatal("failed to create logs dir: " + err.Error())
	}
	logfile, err := os.OpenFile("./logs/requests.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logg.Fatal("failed to open a requests log file: " + err.Error())
	}
	defer logfile.Close()

	// creating router
	router, err := NewRouter(logfile, calendar, config)
	if err != nil {
		logg.Fatal("failed to create a router: " + err.Error())
	}

	// creating server
	server := internalhttp.NewServer(
		config.Server.Host,
		config.Server.Port,
		config.Server.ReadTimeout,
		config.Server.WriteTimeout,
		config.Server.IdleTimeout,
		logg,
		router,
	)

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Fatal("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info(config.App.Name + " is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
