package models

import "time"

type (
	Menu struct {
		ID           uint       `gorm:"primary_key" json:"id"`
		Name         string     `json:"name" gorm:"not null"`
		Description  string     `json:"description"`
		Price        int        `json:"price" gorm:"not null"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
		RestaurantID uint       `json:"restaurantID" gorm:"not null"`
		Restaurant   Restaurant `json:"-"`
	}
)
