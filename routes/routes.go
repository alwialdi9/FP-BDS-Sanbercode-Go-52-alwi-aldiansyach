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

	timeoutval, _ := strconv.Atoi(utils.Getenv("HANDLER_TIMEOUT", "5"))

	r.Use(timeout.Timeout(
		timeout.WithTimeout(time.Duration(timeoutval)*time.Second),
		timeout.WithErrorHttpCode(http.StatusRequestTimeout),                                   // optional
		timeout.WithDefaultMsg(`{"status": "Request Timeout", "msg":"http: Handler timeout"}`), // optional
		timeout.WithCallBack(func(r *http.Request) {
			fmt.Println("timeout happen, url:", r.URL.String())
		}))) // optional

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	MiddlewareRoute := r.Group("/restaurant")
	MiddlewareRoute.Use(middlewares.JwtAuthMiddleware(db))
	MiddlewareRoute.POST("/create", controllers.CreateRestaurant)
	MiddlewareRoute.POST("/create/menus", controllers.CreateMenus)
	// moviesMiddlewareRoute.DELETE("/:id", controllers.DeleteMovie)

	// r.GET("/age-rating-categories", controllers.GetAllRating)
	// r.GET("/age-rating-categories/:id", controllers.GetRatingById)
	// r.GET("/age-rating-categories/:id/movies", controllers.GetMoviesByRatingId)

	// ratingMiddlewareRoute := r.Group("/age-rating-categories")
	// ratingMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
	// ratingMiddlewareRoute.POST("/", controllers.CreateRating)
	// ratingMiddlewareRoute.PATCH("/:id", controllers.UpdateRating)
	// ratingMiddlewareRoute.DELETE("/:id", controllers.DeleteRating)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
