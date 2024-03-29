package controllers

import (
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RestaurantInput struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type MenuParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type MenuRestaurantInput struct {
	ID    string       `json:"restaurant_id"`
	Menus []MenuParams `json:"menus"`
}

type APIRestaurant struct {
	ID   uint
	Name string
	City string
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Get All Restaurant
// @Description Get data for all resto
// @Tags restaurant
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]any
// @Router /get_all_resto [get]
func GetAllRestaurant(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var restaurant []APIRestaurant
	db.Model(&models.Restaurant{}).Find(&restaurant)

	response := []map[string]any{}

	for _, v := range restaurant {
		ratingAvg, totalReview, err := models.SearchReviewByResto(v.ID, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data := map[string]any{"id": v.ID, "name": v.Name, "city": v.City, "total_review": totalReview, "rating_avg": ratingAvg}

		response = append(response, data)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Create data Restaurant
// @Description Create Restaurant
// @Tags restaurant
// @Accept  json
// @Produce  json
// @Param Body body RestaurantInput true "the body to create a restaurant"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param HTTP-X-UID header string true "HTTP-X-UID. Fill with id user"
// @Security BearerToken
// @Success 200 {object} map[string]models.Restaurant
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
	common_req := c.Value("common_request").(models.CommonRequest)

	if !common_req.IsAdmin {
		c.JSON(http.StatusOK, gin.H{"message": "Cannot add restaurant. Please contact admin or SuperAdmin"})
		return
	}

	u.User = common_req.User

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

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Create Menus
// @Description Add menu by id Restaurant
// @Tags Menus
// @Accept  json
// @Produce  json
// @Param Body body MenuRestaurantInput true "the body to create a menu restaurant"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param HTTP-X-UID header string true "HTTP-X-UID. Fill with id user"
// @Security BearerToken
// @Success 200 {object} map[string]string
// @Router /create/menus [post]
func CreateMenus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input MenuRestaurantInput

	common_req := c.Value("common_request").(models.CommonRequest)

	if !common_req.IsAdmin {
		c.JSON(http.StatusOK, gin.H{"message": "Cannot add menu restaurant. Please contact admin or SuperAdmin"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := []models.Menu{}

	restaurant, err := models.SearchRestaurant(input.ID, db)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	for _, v := range input.Menus {
		price, _ := strconv.Atoi(v.Price)
		var data = models.Menu{Name: v.Name, Description: v.Description, Price: price, Restaurant: restaurant, RestaurantID: restaurant.ID}
		u = append(u, data)
	}

	rows, err := models.CreateMenus(db, u)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": `success create ` + strconv.Itoa(rows) + ` Menu restaurant`})
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Delete Menus
// @Description Delete menus by restaurant id
// @Tags menus
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param HTTP-X-UID header string true "HTTP-X-UID. Fill with id user"
// @Security BearerToken
// @Param id path string true "restaurant id"
// @Success 200 {object} map[string]string
// @Router /delete/menus [delete]
func DeleteMenus(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	resto_id := c.Query("restoId")

	if resto_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incomplete Parameter"})
		return
	}

	var menus models.Menu
	if err := db.Where("id = ? AND restaurant_id = ?", c.Param("id"), resto_id).First(&menus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&menus)

	c.JSON(http.StatusOK, gin.H{"data": "Success delete Menus "})
}
