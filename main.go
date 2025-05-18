package main

import (
	"log"
	app "orbital/internal"
	"orbital/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // Optional: don’t fail if missing

	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	server := app.NewServer()

	log.Println("Starting server on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
