package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := UserExistByToken(token); exist {

		_, dbUserInfo := FindUserByToken(token)
		uid := dbUserInfo.Id
		vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
		if err2 != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id format err"})
			return
		}
		err3 := SetFavorite(uid, vid)
		println(FavoriteCount(vid))
		println(IsFavorite(uid, vid))
		if err3 != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: " relate  err"})
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := UserExistByToken(token); exist {

		_, dbUserInfo := FindUserByToken(token)
		uid := dbUserInfo.Id

		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: GetFavoriteList(uid, c.Query("token")),
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

}

//server

func SetFavorite(uid int64, vid int64) error {
	return FavoriteUpdate(uid, vid)
}

// FavoriteList all users have same favorite video list
func GetFavoriteList(uid int64, token string) []Video {
	video_ids, err := FavoriteVid(uid) // GetVideoById(videoId int64) Video
	if err != nil {
		println(video_ids)

	}

	var videos = []Video{
		{
			Id:            1,
			Author:        usersLoginInfo[token],
			PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
			CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
			FavoriteCount: 1,
			CommentCount:  4,
			IsFavorite:    true,
		},
	}
	return videos

}

//dao

func FavoriteVid(uid int64) (vid_list []int64, er error) {

	res := Db.Model(DbFavorite{}).Where("uid=?", uid).Select("vid").Find(&vid_list)
	er = res.Error
	return
}
func FavoriteCount(videoId int64) (count int64) {
	Db.Model(DbFavorite{}).Where("vid=?", videoId).Where("status=?", 1).Select("count(*)").Find(&count)
	return
}
func CommentCount(videoId int64) (count int64) {
	Db.Model(DbComment{}).Where("vid=?", videoId).Select("count(*)").Find(&count)
	return
}
func IsFavorite(userId int64, videoId int64) (re bool) {

	Db.Model(DbFavorite{}).Where("uid=?", userId).Where("vid=?", videoId).Where("status=?", 1).Select("count(*)").Find(&re)

	return

}
func FavoriteUpdate(uid int64, vid int64) error {

	u := DbFavorite{}
	res := Db.Where("uid=?", uid).Where("vid=?", vid).Find(&u)
	if res.RowsAffected == 0 {
		Db.Create(DbFavorite{uid, vid, 1})
	} else if u.Status == 1 {
		Db.Model(DbFavorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", 0)
	} else {
		Db.Model(DbFavorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", 1)
	}

	return nil
}
