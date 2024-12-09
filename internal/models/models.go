package models

import "time"

const (
	StatusAvailable  = "Available"
	StatusDelivering = "Delivering"
	StatusReturning  = "Returning"
	StatusRecharging = "Recharging"
)

type Drone struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	Speed        float64   `json:"speed"`
	Range        float64   `json:"range"`
	ChargingTime int       `json:"chargingTime"`
	Status       string    `json:"status"`
	ReturnTime   time.Time `json:"returnTime"`
}

type Order struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Range     float64   `json:"range"`
	DroneID   string    `json:"droneId"`
	Drone     Drone     `json:"drone" gorm:"foreignKey:DroneID"`
	StartTime time.Time `json:"startTime"`
}
