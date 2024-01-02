package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Restaurant struct {
		ID        uint      `gorm:"primary_key" json:"id"`
		Name      string    `json:"name"`
		City      string    `json:"city"`
		Reviews   []Review  `json:"-"`
		Menus     []Menu    `json:"-"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		UserID    uint      `json:"user_id" gorm:"not null"`
		User      User      `json:"-"`
	}
)

func (u *Restaurant) SaveRestaurant(db *gorm.DB) (*Restaurant, error) {

	var err error = db.Create(&u).Error
	if err != nil {
		return &Restaurant{}, err
	}
	return u, nil
}

func SearchRestaurant(id string, db *gorm.DB) (Restaurant, error) {
	var err error
	resto := Restaurant{}
	err = db.Model(resto).Where("id = ?", id).Take(&resto).Error
	if err != nil {
		return resto, err
	}
	return resto, nil
}
