package controllers

import (
	"final-project/models"
	"final-project/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

type ResetLinkInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPassInput struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
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

	token, username, id, err := models.LoginCheck(u.Email, u.Password, db)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password incorrect"})
		return
	}

	user := map[string]string{
		"username": username,
		"email":    u.Email,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login Success", "user": user, "token": token, "id": id})
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

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Get Reset Link
// @Description get Link by ID
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Body body ResetLinkInput true "the body to get reset link password"
// @Success 200 {object} map[string]any
// @Router /get_reset_link [post]
func GetResetLink(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input ResetLinkInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := models.User{}
	err := db.Model(models.User{}).Where("email = ?", input.Email).Take(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := u.Email

	if email == "" {
		c.JSON(http.StatusOK, gin.H{"message": "User not found"})
		return
	}

	resetToken := utils.GetToken(10)
	u.ResetCode = resetToken
	u.ResetTime = time.Now().Add(time.Minute * 15)

	db.Save(&u)

	response := map[string]any{
		"email":        u.Email,
		"reset_link":   c.Request.Host + "/reset_password/" + u.ResetCode,
		"link_expired": u.ResetTime.Format("02 January 2006 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Reset Password
// @Description reset password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Body body ResetPassInput true "the body to reset password"
// @Success 200 {object} map[string]string
// @Router /reset_password [post]
func ResetPassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input ResetPassInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := c.Param("token")

	u := models.User{}
	err := db.Model(models.User{}).Where("reset_code = ?", token).Take(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resetTime := u.ResetTime

	if resetTime.Before(time.Now()) {
		c.JSON(http.StatusOK, gin.H{"message": "link expired, please call reset_link endpoint again"})
		return
	}

	password := input.Password
	new_password := input.NewPassword

	errPass := models.VerifyPassword(password, u.Password)

	if errPass != nil && errPass == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": errPass.Error()})
		return
	}

	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if errPassword != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errPassword.Error()})
		return
	}
	u.Password = string(hashedPassword)
	u.ResetCode = ""
	db.Save(&u)

	response := map[string]any{
		"message": "Success reset password",
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}
