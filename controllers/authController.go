package controllers

import (
	"final-project/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	admin = "admin"
	user  = "user"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Login User
// @Description Logging to get jwt token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Body body LoginInput true "the body to login a user"
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := models.User{}

	u.Email = input.Email
	u.Password = input.Password

	token, username, err := models.LoginCheck(u.Email, u.Password, db)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password incorrect"})
		return
	}

	user := map[string]string{
		"username": username,
		"email":    u.Email,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login Success", "user": user, "token": token})
}

// Register godoc
// @Summary Register a user.
// @Description registering a user from public access.
// @Tags Auth
// @Param Body body RegisterInput true "the body to register a user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Email = input.Email
	u.Password = input.Password
	u.Role = input.Role
	if input.Role == "" {
		u.Role = user
	} else if strings.ToLower(input.Role) == "admin" {
		u.Role = admin
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "registration failed", "error": "Role for account only user or admin"})
		return
	}

	_, err := u.SaveUser(db)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := map[string]string{
		"username": input.Username,
		"email":    input.Email,
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success", "user": user})

}
