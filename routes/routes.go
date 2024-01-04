package routes

import (
	"final-project/controllers"
	"final-project/middlewares"
	"final-project/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	// set db to gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	gin.SetMode(gin.ReleaseMode)

	timeoutval, _ := strconv.Atoi(utils.Getenv("HANDLER_TIMEOUT", "5"))

	r.Use(timeout.Timeout(
		timeout.WithTimeout(time.Duration(timeoutval)*time.Second),
		timeout.WithErrorHttpCode(http.StatusRequestTimeout),
		timeout.WithDefaultMsg(`{"status": "Request Timeout", "msg":"http: Handler timeout"}`),
		timeout.WithCallBack(func(r *http.Request) {
			fmt.Println("timeout happen, url:", r.URL.String())
		})))

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/get_reset_link", controllers.GetResetLink)
	r.POST("/reset_password/:token", controllers.ResetPassword)

	MiddlewareRoute := r.Group("/restaurant")
	MiddlewareRoute.Use(middlewares.JwtAuthMiddleware(db))
	MiddlewareRoute.POST("/create", controllers.CreateRestaurant)
	MiddlewareRoute.POST("/create/menus", controllers.CreateMenus)
	MiddlewareRoute.DELETE("/delete/menus/:id", controllers.DeleteMenus)

	r.GET("/get_all_resto", controllers.GetAllRestaurant)

	UserMiddlewareRoute := r.Group("/user")
	UserMiddlewareRoute.Use(middlewares.JwtAuthMiddleware(db))
	UserMiddlewareRoute.POST("/create/orders", controllers.CreateOrder)
	UserMiddlewareRoute.GET("/show/order/:id/restaurant", controllers.ShowOrderByResto)
	UserMiddlewareRoute.POST("/send_review", controllers.CreateReview)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
