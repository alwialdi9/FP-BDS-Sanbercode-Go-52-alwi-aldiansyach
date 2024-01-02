package controllers

import (
	"final-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RestaurantInput struct {
	Name string `json:"name"`
	City string `json:"city"`
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Create data Restaurant
// @Description Create Restaurant
// @Tags restaurant
// @Accept  json
// @Produce  json
// @Param Body body RestaurantInput true "the body to create a restaurant"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]interface{}
// @Router /restaurant/create [post]
func CreateRestaurant(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input RestaurantInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := models.Restaurant{}

	u.Name = input.Name
	u.City = input.City
	user := models.CheckAdmin(c.Request.Header.Get("HTTP-X-UID"), db)

	if user.Role != "admin" {
		c.JSON(http.StatusOK, gin.H{"message": "Cannot add restaurant. Please contact admin or SuperAdmin"})
		return
	}

	u.User = user

	_, err := u.SaveRestaurant(db)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant := map[string]string{
		"name": input.Name,
		"city": input.City,
	}

	c.JSON(http.StatusOK, gin.H{"message": "success create restaurant", "data": restaurant})
}
