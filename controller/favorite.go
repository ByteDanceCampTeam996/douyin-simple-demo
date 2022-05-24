package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		uid := usersLoginInfo[token].Id
		vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
		if err2 != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id format err"})
			return
		}
		err3 := SetFavorite(uid, vid)
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
	uid := usersLoginInfo[token].Id

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: GetFavoriteList(uid, c.Query("token")),
	})

}

//server

func SetFavorite(uid int64, vid int64) error {
	return FavoriteUpdate(uid, vid)
}

func GetFavoriteList(uid int64, token string) []Video {
	video_ids, err := FavoriteVid(uid)
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

type Favorite struct {
	Uid    int64
	Vid    int64
	Status int
}

func FavoriteVid(uid int64) (vid_list []int64, er error) {

	//db.LogMode(true)
	res := Db.Table("favorites").Where("uid=?", uid).Select("vid").Find(&vid_list)
	er = res.Error
	return
}
func FavoriteUpdate(uid int64, vid int64) error {

	u := Favorite{}
	res :=Db.Where("uid=?", uid).Where("vid=?", vid).Find(&u)
	if res.RowsAffected == 0 {
		Db.Create(Favorite{uid, vid, 1})
	} else if u.Status == 1 {
		Db.Model(Favorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", 0)
	} else {
		Db.Model(Favorite{}).Where("uid=?", uid).Where("vid=?", vid).Update("status", 1)
	}

	return nil
}
