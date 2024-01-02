package models

import "time"

type (
	OrderHistory struct {
		ID           uint       `json:"id" gorm:"primary_key"`
		TotalPrice   int        `json:"total_price" gorm:"not null"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
		RestaurantID uint       `json:"restaurantID" gorm:"not null"`
		Restaurant   Restaurant `json:"-"`
		UserID       uint       `json:"userID"`
		User         User       `json:"-"`
		Menus        []Menu     `json:"-" gorm:"many2many:menu;"`
	}
)
