package main

import (
	"log"
	"spageti-x-drone-management/internal/config"
	"spageti-x-drone-management/internal/database"
	"spageti-x-drone-management/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	// Initialize handlers with database
	h := handlers.NewHandler(db)

	// Route definitions
	r.POST("/drones", h.AddDrone)
	r.GET("/drones", h.GetDrones)
	r.PUT("/drones/:id", h.UpdateDrone)
	r.DELETE("/drones/:id", h.RemoveDrone)

	r.POST("/orders", h.CreateOrder)
	r.GET("/orders", h.GetOrders)

	// Start background worker
	go h.UpdateDroneStatuses()

	// Start server
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
