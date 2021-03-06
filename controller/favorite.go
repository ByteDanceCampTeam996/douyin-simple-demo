package controller

import (
	"net/http"
	"strconv"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/service"
	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := dao.UserExistByToken(token); exist {

		_, dbUserInfo := dao.FindUserByToken(token)
		uid := dbUserInfo.Id
		vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
		if err2 != nil {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "video_id format err"})
			return
		}
		action_type, _ := strconv.Atoi(c.Query("action_type"))
		err3 := service.SetFavorite(uid, vid, action_type)
		if err3 != nil {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: " relate  err"})
			return
		}

		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := dao.UserExistByToken(token); !exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	_, dbUserInfo := dao.FindUserByToken(token)
	uid := dbUserInfo.Id

	err, videoList := service.GetFavoriteList(uid)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  err.Error(),
			},
			VideoList: videoList,
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
