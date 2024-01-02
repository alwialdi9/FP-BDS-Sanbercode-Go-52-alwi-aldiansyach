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
		common_req := models.CommonRequest{}
		err = db.Model(u).Where("id = ?", user_id).Take(&u).Error
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		if u.Role == "admin" {
			common_req.IsAdmin = true
		}
		common_req.User = u
		c.Set("common_request", common_req)
		c.Next()
	}
}
