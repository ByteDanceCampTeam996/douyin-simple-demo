package controller

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDb initial database for local test
func InitDb() {
	// 设置连接MySQL数据库的账号、密码，以及连接的数据库
	//dsn := "root:qwer1234@tcp(127.0.0.1:3306)/userlogininfo?charset=utf8mb4&parseTime=True&loc=Local"
       dsn := "root:123456@tcp(127.0.0.1:3306)/douyin_db?charset=utf8&parseTime=True&loc=Local"
	db1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db1
	db.AutoMigrate(&DbUserInfo{})
	db.Model(&DbUserInfo{}).Count(&userIdSequence)
	fmt.Println(userIdSequence)
}
