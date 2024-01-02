package models

import (
	"final-project/utils/token"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	User struct {
		ID           uint           `json:"id" gorm:"primary_key"`
		Username     string         `gorm:"not null;unique" json:"username"`
		Email        string         `json:"email" gorm:"not null;unique"`
		Password     string         `json:"password"`
		Role         string         `json:"role"`
		CreatedAt    time.Time      `json:"created_at"`
		UpdatedAt    time.Time      `json:"updated_at"`
		Review       []Review       `json:"-"`
		OrderHistory []OrderHistory `json:"-"`
	}
)

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CheckAdmin(id string, db *gorm.DB) User {
	var err error
	u := User{}
	err = db.Model(User{}).Where("id = ?", id).Take(&u).Error
	if err != nil {
		return User{}
	}
	return u
}

func LoginCheck(email string, password string, db *gorm.DB) (string, string, error) {
	var err error

	u := User{}

	err = db.Model(User{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		return "", "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", "", err
	}
	return token, u.Username, nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	//turn password into hash
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		return &User{}, errPassword
	}
	u.Password = string(hashedPassword)
	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	var err error = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}
