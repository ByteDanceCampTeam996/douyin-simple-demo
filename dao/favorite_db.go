package dao

import (
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
)

func FavoriteVid(uid int64) (videoList []int64, er error) {
	res := Db.Model(model.DbFavorite{}).Where("uid=?", uid).Select("vid").Find(&videoList)
	er = res.Error
	return
}
func FavoriteCount(videoId int64) (count int64) {
	Db.Model(model.DbFavorite{}).Where("vid=?", videoId).Where("status=?", 1).Select("count(*)").Find(&count)
	return
}
func CommentCount(videoId int64) (count int64) {
	Db.Model(model.DbComment{}).Where("vid=?", videoId).Select("count(*)").Find(&count)
	return
}
func IsFavorite(userId int64, videoId int64) (re bool) {

	Db.Model(model.DbFavorite{}).Where("uid=?", userId).Where("vid=?", videoId).Where("status=?", 1).Select("count(*)").Find(&re)

	return

}
func FavoriteUpdate(uid int64, vid int64) error {

	u := model.DbFavorite{}
	res := Db.Where("uid=?", uid).Where("vid=?", vid).Find(&u)
	if res.RowsAffected == 0 {
		Db.Create(model.DbFavorite{Uid: uid, Vid: vid, Status: 1})
	} else if u.Status == 1 {
		Db.Model(model.DbFavorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", 0)
	} else {
		Db.Model(model.DbFavorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", 1)
	}

	return nil
}
