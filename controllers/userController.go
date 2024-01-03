package controllers

import (
	"final-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Menu struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderMenuInput struct {
	RestaurantID string `json:"restaurant_id"`
	Menus        []Menu `json:"order_menu"`
}

type ResponseOrder struct {
	ID         uint   `json:"id"`
	TotalPrice int    `json:"total_price"`
	Menus      []Menu `json:"-" gorm:"many2many:orderHistory_menus;"`
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Create order
// @Description Create order user by resto
// @Tags User
// @Accept  json
// @Produce  json
// @Param Body path string true "Account ID"
// @Success 200 {object} map[string]models.OrderHistory
// @Router /create/orders [post]
func CreateOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input OrderMenuInput

	user := c.Value("common_request").(models.CommonRequest).User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	menus := []models.Menu{}
	orders := models.OrderHistory{}
	resto_id := input.RestaurantID

	if resto_id == "nil" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant not found"})
		return
	}
	restaurant, err := models.SearchRestaurant(resto_id, db)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orders.User = user
	orders.Restaurant = restaurant

	var price int
	for _, v := range input.Menus {
		menu, total, err := models.CountTotalPrice(v.ID, v.Quantity, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		price += total
		menus = append(menus, menu)
	}
	orders.TotalPrice = price
	orders.Menus = menus

	orderMenus, err := models.CreateOrders(db, orders)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": orderMenus})
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Show Order By Resto
// @Description Get Order by id
// @Tags Order User
// @Accept  json
// @Produce  json
// @Param id path string true "Restaurant id"
// @Success 200 {object} models.OrderHistory
// @Router /show/order/:id/restaurant [get]
func ShowOrderByResto(c *gin.Context) {
	var orders []models.OrderHistory

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Model(&models.OrderHistory{}).Preload("Menus").Where("restaurant_id = ?", c.Param("id")).Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var result []map[string]any

	for _, v := range orders {
		data := map[string]interface{}{"id": v.ID, "restaurant_id": v.RestaurantID, "total_price": v.TotalPrice, "order": v.Menus}

		result = append(result, data)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}