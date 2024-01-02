package routes

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
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

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	moviesMiddlewareRoute := r.Group("/restaurant")
	moviesMiddlewareRoute.Use(middlewares.JwtAuthMiddleware(db))
	moviesMiddlewareRoute.POST("/create", controllers.CreateRestaurant)
	// moviesMiddlewareRoute.POST("/order/history", controllers.UpdateMovie)
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
