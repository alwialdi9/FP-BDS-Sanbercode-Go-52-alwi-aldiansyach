package models

import (
	"time"

	"gorm.io/gorm"
)

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
		Menus        []Menu     `json:"-" gorm:"many2many:orderHistory_menus;"`
	}
)

func CreateOrders(db *gorm.DB, u OrderHistory) (OrderHistory, error) {

	orders := u
	var err error = db.Create(&orders).Error

	return orders, err
}
