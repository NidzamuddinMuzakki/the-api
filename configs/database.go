package configs

import (
	"fmt"
	"os"

	"github.com/NidzamuddinMuzakki/the-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect() error {
	var err error
	MYSQL_HOST := os.Getenv("MYSQL_HOST")
	MYSQL_DB := os.Getenv("MYSQL_DATABASE")
	MYSQL_USER := os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")

	MYSQL_HOST = "localhost:3306"
	MYSQL_DB = "the-api"
	MYSQL_USER = "root"
	MYSQL_PASSWORD = "123"

	DATABASE_URI := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DB)
	Database, err = gorm.Open(mysql.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic(err)
	}

	Database.AutoMigrate(&models.User{}, &models.Article{})

	return nil
}
