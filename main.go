package main

import (
	"log"
	"os"
	"simple-shortener/config"
	"simple-shortener/handler"
	"simple-shortener/repository"
	"simple-shortener/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	dbPath := cfg.DatabasePath
	if p := os.Getenv("SHORTENER_DB"); p != "" {
		dbPath = p
	}

	store, err := repository.NewGormStore(dbPath)
	if err != nil {
		log.Fatalf("failed init store: %v", err)
	}

	svc := service.NewService(store)

	// Pass baseURL to handler
	h := handler.NewHandler(svc, cfg.BaseURL)

	// Start HTTP server (Gin)
	log.Printf("Starting server at %s", cfg.ServerPort)
	log.Printf("Base URL: %s", cfg.BaseURL)
	if err := h.Run(cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
