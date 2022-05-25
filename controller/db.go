package controller

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDb initial database for local test
var Db *gorm.DB

func ConnectDB() {
	var (
		err error
	)
	user := "root"
	password := "123456"
	host := "127.0.0.1:3306"
	dbname := "douyin"
	//dsn := "root:123456@tcp(127.0.0.1:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user,
		password,
		host,
		dbname)
	Db, err = gorm.Open(mysql.Open(connectStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//自动生成表结构
	dberr := Db.AutoMigrate(&DbUserInfo{})
	if dberr != nil {
		println(err)
	}
	dberr = Db.AutoMigrate(&UserFollowInfo{})
	if dberr != nil {
		println(err)
	}
	dberr = Db.AutoMigrate(&Follow{})
	if dberr != nil {
		println(err)
	}
	dberr = Db.AutoMigrate(&DbFavorite{})
	if dberr != nil {
		println(err)
	}
	dberr = Db.AutoMigrate(&DbComment{})
	if dberr != nil {
		println(err)
	}
	//读取数据库中现有的用户数量
	Db.Model(&DbUserInfo{}).Count(&userIdSequence)
	fmt.Println(userIdSequence)

}
