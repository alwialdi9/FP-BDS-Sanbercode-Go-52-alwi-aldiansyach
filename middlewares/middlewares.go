package middlewares

import (
	"final-project/models"
	"final-project/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JwtAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		user_id := c.Request.Header.Get("HTTP-X-UID")
		if err != nil || user_id == "" {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		u := models.User{}
		err = db.Model(u).Where("id = ?", user_id).Take(&u).Error
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
