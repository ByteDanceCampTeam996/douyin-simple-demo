package dao

import (
	"fmt"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
)

// GetVideoById 根据视频ID获取视频完整信息
func GetVideoById(curUserId int64, videoId int64) (error, model.Video) {
	var video model.Video
	var dbVideoInfo model.DbVideoInfo
	// 视频信息
	if err := Db.Find(&dbVideoInfo, "video_id = ?", videoId).Error; err != nil {
		fmt.Println("视频信息查询失败！")
		return err, video
	}
	// 作者信息
	var dbUserInfo model.DbUserInfo
	if err := Db.Find(&dbUserInfo).Where("user_id = ?", dbVideoInfo.UserId).Error; err != nil {
		fmt.Println("视频作者信息查询失败！")
		return err, video
	}
	// 获取作者关注数和被关注数
	var follow model.Follow
	var followCount int64
	var fansCount int64
	if err := Db.Where("follow_id = ? AND status <> 1", dbVideoInfo.UserId).Find(&follow).Count(&fansCount).Error; err != nil {
		fmt.Println("查不到视频作者粉丝数！")
		return err, video
	}
	if err := Db.Where("user_id = ? AND status <> 1", dbVideoInfo.UserId).Find(&follow).Count(&followCount).Error; err != nil {
		fmt.Println("查不到视频作者关注数！")
		return err, video
	}
	// To do: 添加当前用户id字段后判断
	// 当前用户是否关注作者，传的参数字段不够判断
	var isFollow bool
	isFollow = false
	// 当前用户是否点赞了视频,需要传当前登陆的id才能查
	var isFavorite bool
	isFavorite = IsFavorite(curUserId, videoId)

	// 点赞数
	var favoriteCount int64
	favoriteCount = FavoriteCount(videoId)
	fmt.Printf("\n视频点赞数：%d", favoriteCount)
	// 评论数
	var commentCount int64
	commentCount = CommentCount(videoId)
	fmt.Printf("\n视频评论数：%d", commentCount)
	// 拼接返回结果
	var author model.User
	author = model.User{Id: dbVideoInfo.UserId, Name: dbUserInfo.UserName, FollowCount: followCount, FollowerCount: fansCount, IsFollow: isFollow}
	video = model.Video{Id: videoId, PlayUrl: dbVideoInfo.PlayUrl, CoverUrl: dbVideoInfo.CoverUrl, Title: dbVideoInfo.Title,
		FavoriteCount: favoriteCount, CommentCount: commentCount, IsFavorite: isFavorite, Author: author}
	// To do: 更多异常处理考虑
	return nil, video
}

func FavoriteVid(uid int64) (vid_list []int64, er error) {

	res := Db.Model(model.DbFavorite{}).Where("uid=?", uid).Select("vid").Find(&vid_list)
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
		Db.Create(model.DbFavorite{uid, vid, 1})
	} else if u.Status == 1 {
		Db.Model(model.DbFavorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", 0)
	} else {
		Db.Model(model.DbFavorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", 1)
	}

	return nil
}
