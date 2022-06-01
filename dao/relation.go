package dao

import "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"

//获取关注列表里面的一条记录
func GetFollowRecord(userId, followId int64) (follow model.Follow, rows int64) {
	res := Db.Where("user_id = ? AND follow_id = ?", userId, followId).Find(&follow)
	return follow, res.RowsAffected
}

//更新关注表里面的关注状态
func UpdateFollowStatus(userId, followId, status int64) error {
	res := Db.Model(&model.Follow{}).Where("user_id = ? AND follow_id = ?",
		userId, followId).Update("status", status)
	return res.Error
}

//在关注表里创建一条关注信息
func CreateFollow(userId, followId, status int64) error {
	res := Db.Model(&model.Follow{}).Create(&model.Follow{UserId: userId, FollowId: followId, Status: status})
	return res.Error
}

//更新用户粉丝关注数量表
func UpdateUserFollowInfo(userId, follows, fans int64) error {
	res := Db.Model(&model.UserFollowInfo{}).Where("user_id = ?", userId).Select("follow_count", "follower_count").
		Updates(&model.UserFollowInfo{FollowCount: follows, FollowerCount: fans})
	return res.Error
}

//获取用户的粉丝数和关注数
func GetUserFollowInfo(userId int64) (userinfo model.UserFollowInfo) {
	Db.Where("user_id = ?", userId).Find(&userinfo)
	return
}

//获取用户的所有关注对象
func GetFollows(userId int64) (follow []model.Follow, err error) {
	res := Db.Where("user_id = ? AND status <> 0", userId).Find(&follow)
	return follow, res.Error
}

//获取用户的所有粉丝
func GetFans(userId int64) (follow []model.Follow, err error) {
	res := Db.Where("follow_id = ? AND status <> 0", userId).Find(&follow)
	return follow, res.Error
}
