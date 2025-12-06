package main

import (
	"log"
	"os"
	"simple-shortener/handler"
	repostitory "simple-shortener/repository"
	"simple-shortener/service"
)

func main() {
	dbPath := "shortener.db"
	if p := os.Getenv("SHORTENER_DB"); p != "" {
		dbPath = p
	}

	store, err := repostitory.NewGormStore(dbPath)
	if err != nil {
		log.Fatalf("failed init store: %v", err)
	}

	svc := service.NewService(store)

	h := handler.NewHandler(svc)

	// Start HTTP server (Gin)
	if err := h.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
