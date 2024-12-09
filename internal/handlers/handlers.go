package handlers

import (
	"net/http"
	"time"

	"spageti-x-drone-management/internal/database"
	"spageti-x-drone-management/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	db *database.Database
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) AddDrone(c *gin.Context) {
	var drone models.Drone
	if err := c.ShouldBindJSON(&drone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drone.ID = uuid.New().String()
	drone.Status = models.StatusAvailable

	if err := h.db.DB.Create(&drone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create drone"})
		return
	}

	c.JSON(http.StatusCreated, drone)
}

func (h *Handler) GetDrones(c *gin.Context) {
	var drones []models.Drone
	status := c.Query("status")
	droneRange := c.Query("range")

	if status != "" {
		h.db.DB.Where("status = ?", status)
	}

	if droneRange != "" {
		h.db.DB.Where("range = ?", droneRange)
	}

	if err := h.db.DB.Find(&drones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drones"})
		return
	}

	c.JSON(http.StatusOK, drones)
}

func (h *Handler) UpdateDrone(c *gin.Context) {
	id := c.Param("id")
	var drone models.Drone

	if err := h.db.DB.First(&drone, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Drone not found"})
		return
	}

	if err := c.ShouldBindJSON(&drone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.DB.Save(&drone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update drone"})
		return
	}

	c.JSON(http.StatusOK, drone)
}

func (h *Handler) RemoveDrone(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.DB.Delete(&models.Drone{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete drone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drone deleted successfully"})
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var drone models.Drone
	if err := h.db.DB.First(&drone, "status = ?", models.StatusAvailable).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No available drones"})
		return
	}

	order.ID = uuid.New().String()
	order.DroneID = drone.ID
	order.StartTime = time.Now()

	drone.Status = models.StatusDelivering
	drone.ReturnTime = time.Now().Add(time.Hour)

	tx := h.db.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	if err := tx.Save(&drone).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update drone status"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, order)
}

func (h *Handler) GetOrders(c *gin.Context) {
	var orders []models.Order
	if err := h.db.DB.Preload("Drone").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) UpdateDroneStatuses() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		var drones []models.Drone
		h.db.DB.Where("status IN ?", []string{models.StatusDelivering, models.StatusReturning}).Find(&drones)

		for _, drone := range drones {
			if time.Now().After(drone.ReturnTime) {
				drone.Status = models.StatusRecharging
				h.db.DB.Save(&drone)
			}
		}
	}
}
