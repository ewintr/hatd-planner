package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	apiPort    = flag.String("port", "8092", "api port")
	apiKey     = flag.String("key", "testKey", "api key")
	dbHost     = flag.String("dbhost", "localhost", "database host")
	dbPort     = flag.String("dbport", "5432", "database port")
	dbName     = flag.String("dbname", "planner", "database name")
	dbUser     = flag.String("dbuser", "test", "database user")
	dbPassword = flag.String("dbpassword", "test", "database password")
)

func main() {
	flag.Parse()

	repo, err := NewPostgres(*dbHost, *dbPort, *dbName, *dbUser, *dbPassword)
	if err != nil {
		fmt.Printf("could not open postgres db: %s", err.Error())
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("configuration", "configuration", map[string]string{
		"port":   *apiPort,
		"dbHost": *dbHost,
		"dbPort": *dbPort,
		"dbName": *dbName,
		"dbUser": *dbUser,
	})
	recurrer := NewRecur(repo, repo, logger)
	go recurrer.Run(12 * time.Hour)

	srv := NewServer(repo, *apiKey, logger)
	go http.ListenAndServe(fmt.Sprintf(":%s", *apiPort), srv)

	logger.Info("service started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("service stopped")
}
