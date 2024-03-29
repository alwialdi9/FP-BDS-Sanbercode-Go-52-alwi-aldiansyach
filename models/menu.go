package models

import (
	"time"

	"gorm.io/gorm"
)

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

func CreateMenus(db *gorm.DB, u []Menu) (int, error) {

	menus := u
	result := db.Create(menus)

	return int(result.RowsAffected), result.Error
}

func CountTotalPrice(id string, quantity int, db *gorm.DB) (Menu, int, error) {
	var err error
	menu := Menu{}
	err = db.Where("id = ?", id).Take(&menu).Error
	if err != nil {
		return menu, 0, err
	}
	price := menu.Price * quantity
	return menu, price, nil
}
