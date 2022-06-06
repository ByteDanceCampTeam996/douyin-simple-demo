package dao

import (
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
)

//获取点赞视频id列表
func FavoriteVid(uid int64) (videoList []int64, er error) {
	res := Db.Model(model.DbFavorite{}).Where("uid=?", uid).Select("vid").Where("status=?", 1).Find(&videoList)
	er = res.Error
	return
}

//获取视频点赞数量
func FavoriteCount(videoId int64) (count int64) {
	Db.Model(model.DbFavorite{}).Where("vid=?", videoId).Where("status=?", 1).Select("count(*)").Find(&count)
	return
}

//获取视频评论数量
func CommentCount(videoId int64) (count int64) {
	Db.Model(model.DbComment{}).Where("vid=?", videoId).Select("count(*)").Find(&count)
	return
}

//判断视频是否点赞
func IsFavorite(userId int64, videoId int64) (re bool) {

	Db.Model(model.DbFavorite{}).Where("uid=?", userId).Where("vid=?", videoId).Where("status=?", 1).Select("count(*)").Find(&re)

	return

}

//更新点赞状态
func FavoriteUpdate(uid int64, vid int64, status int) error {

	u := model.DbFavorite{}
	res := Db.Where("uid=?", uid).Where("vid=?", vid).Find(&u)
	if res.RowsAffected == 0 {
		Db.Create(model.DbFavorite{Uid: uid, Vid: vid, Status: status})
	} else {
		Db.Model(model.DbFavorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", status)
	}

	return nil
}
