package main

import (
	"log"
	"todo-list/backend/config"
	"todo-list/backend/database"
	"todo-list/backend/router"
)

func main() {
	cfg := config.Load()

	if err := database.Init(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized")

	r := router.Setup(cfg.CorsOrigin)

	log.Printf("Server starting on :%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
