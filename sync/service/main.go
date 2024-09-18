package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	port := os.Getenv("PLANNER_PORT")
	apiKey := os.Getenv("PLANNER_API_KEY")
	if apiKey == "" {
		fmt.Println("PLANNER_API_KEY is empty")
		os.Exit(1)
	}

	dbHost := os.Getenv("PLANNER_DB_HOST")
	dbPort := os.Getenv("PLANNER_DB_PORT")
	dbName := os.Getenv("PLANNER_DB_NAME")
	dbUser := os.Getenv("PLANNER_DB_USER")
	dbPassword := os.Getenv("PLANNER_DB_PASSWORD")
	repo, err := NewPostgres(dbHost, dbPort, dbName, dbUser, dbPassword)
	if err != nil {
		fmt.Printf("could not open sqlite db: %s", err.Error())
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("configuration", "configuration", map[string]string{
		"port":   port,
		"dbHost": dbHost,
		"dbPort": dbPort,
		"dbName": dbName,
		"dbUser": dbUser,
	})

	srv := NewServer(repo, apiKey, logger)
	go http.ListenAndServe(fmt.Sprintf(":%s", port), srv)

	logger.Info("service started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("service stopped")
}
