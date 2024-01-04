package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

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

func (review *Review) CreateReview(db *gorm.DB) (*Review, error) {

	var err error = db.Create(&review).Error
	if err != nil {
		return &Review{}, err
	}
	return review, nil
}

func SearchReviewByResto(id uint, db *gorm.DB) (float64, int, error) {
	var err error
	var result struct {
		RatingAvg   float64 `json:"rating_avg"`
		TotalReview int     `json:"total_review"`
	}
	err = db.Model(&Review{}).Select("AVG(rating) as rating_avg, COUNT(1) as total_review").Where("restaurant_id = ?", id).Take(&result).Error
	if err != nil {
		return 0, 0, err
	}
	fmt.Println(result)
	if result.TotalReview > 0 {
		return result.RatingAvg, result.TotalReview, nil
	}
	return 0, result.TotalReview, nil
}
