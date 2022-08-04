package main

import (
	"flag"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"
	"os"

	_ "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/migrations"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	flags      = flag.NewFlagSet("goose", flag.ExitOnError)
	dir        = flags.String("dir", "migrations", "directory with migration files")
	configFile = flags.String("config", "configs/.env", "path to configuration file")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("failed to parse flags: %v\n", err)
	}

	args := flags.Args()
	fmt.Println(args)

	if err := godotenv.Load(*configFile); err != nil {
		log.Fatalf("failed to load .env file: %v\n", err)
	}

	config := NewConfig()
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to populate config: %v\n", err)
	}

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	dsn := config.PostgreSQL.BuildDSN()

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	if command == "create" {
		tmp := []string{"pgx", dsn}
		args = append(tmp, args...)
		args = append(args, "go")
	}

	arguments := make([]string, 0)
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err) //nolint:gocritic
	}
}
