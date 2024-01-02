package config

import (
	"final-project/models"
	"final-project/utils"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDataBase() *gorm.DB {
	username := utils.Getenv("MYSQLUSER", "root")
	password := utils.Getenv("MYSQL_ROOT_PASSWORD", "")
	host := utils.Getenv("MYSQLHOST", "tcp(127.0.0.1:3306)")
	database := utils.Getenv("MYSQL_DATABASE", "restaurant_api")

	dsn := fmt.Sprintf("%v:%v@%v/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.Restaurant{}, &models.Review{}, &models.User{}, &models.Menu{}, &models.OrderHistory{})

	return db
}
