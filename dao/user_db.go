package dao

import (
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
)

// 通过Token判断用户是否存在
func UserExistByToken(token string) (error, bool) {
	var dbUserInfo DbUserInfo
	if result := Db.Where("token = ?", token).First(&dbUserInfo); result.Error == nil {
		return nil, true
	} else {
		return result.Error, false
	}
}

// 通过用户名判断用户是否存在
func UserExistByUsername(username string) (error, bool) {
	var dbUserInfo DbUserInfo
	if result := Db.Where("user_name = ?", username).First(&dbUserInfo); result.Error == nil {
		return nil, true
	} else {
		return result.Error, false
	}
}

// 通过Token返回用户基本登录信息
func FindUserByToken(token string) (error, DbUserInfo) {
	var dbUserInfo DbUserInfo
	result := Db.Where("token = ?", token).First(&dbUserInfo)
	return result.Error, dbUserInfo

}

// 通过用户名返回用户基本登录信息
func FindUserByUsername(username string) (error, DbUserInfo) {
	var dbUserInfo DbUserInfo
	result := Db.Where("user_name = ?", username).First(&dbUserInfo)
	return result.Error, dbUserInfo

}

func FindFollowInfoById(id int64) (error, UserFollowInfo) {
	var userFollowInfo UserFollowInfo
	result := Db.Where("user_id = ?", id).First(&userFollowInfo)
	return result.Error, userFollowInfo

}

func CreateNewUser(newUser DbUserInfo) {
	Db.Create(newUser)
}

func CreateNewFollowInfo(newFollowInfo UserFollowInfo) {
	Db.Create(newFollowInfo)
}
