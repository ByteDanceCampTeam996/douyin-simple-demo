package dao

import (
	"fmt"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"

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
	dbErr := Db.AutoMigrate(&model.DbUserInfo{})
	if dbErr != nil {
		println(err)
	}

	dbErr = Db.AutoMigrate(&model.DbVideoInfo{})
	if dbErr != nil {
		println(err)
	}
	dbErr = Db.AutoMigrate(&model.Follow{})
	if dbErr != nil {
		println(err)
	}
	dbErr = Db.AutoMigrate(&model.DbFavorite{})
	if dbErr != nil {
		println(err)
	}
	dbErr = Db.AutoMigrate(&model.DbComment{})
	if dbErr != nil {
		println(err)
	}
	dbErr = Db.AutoMigrate(&model.UserFollowInfo{})
	if dbErr != nil {
		println(err)
	}
	/*//读取数据库中现有的用户数量
	Db.Model(&model.DbUserInfo{}).Count(&model.userIdSequence)
	fmt.Println(controller.userIdSequence)*/

}
