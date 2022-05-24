package controller

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDb initial database for local test
func InitDb() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	db1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db1
	db.AutoMigrate(&DbFavorite{})
	db.Model(&DbFavorite{}).Count(&userIdSequence)
	db.AutoMigrate(&DbComment{})
	db.Model(&DbComment{}).Count(&userIdSequence)
	fmt.Println(userIdSequence)
}
