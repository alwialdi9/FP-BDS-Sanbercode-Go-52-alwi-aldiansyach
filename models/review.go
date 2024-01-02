package models

import "time"

type (
	Review struct {
		ID           uint       `gorm:"primary_key" json:"id"`
		Rating       int        `json:"rating" gorm:"not null"`
		Content      string     `json:"content"`
		UserID       uint       `json:"userID" gorm:"not null"`
		User         User       `json:"-"`
		RestaurantID uint       `json:"restaurantID" gorm:"not null"`
		Restaurant   Restaurant `json:"-"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
	}
)
