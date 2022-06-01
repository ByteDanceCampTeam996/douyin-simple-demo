package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/service"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed 首页视频列表获取接口
func Feed(c *gin.Context) {
	// 参数获取与校验
	latestTime := c.Query("latest_time")
	token := c.Query("token")
	var videoList []model.Video
	var nextTime int64
	var err error
	// 字符串字段格式处理
	var nextVideoTime int64
	if nextVideoTime, err = strconv.ParseInt(latestTime, 10, 64); err != nil {
		fmt.Println("latest_time字符串转换为int64失败！")
		nextVideoTime = 0
	}
	if err, videoList, nextTime = service.GetVideoList(token, nextVideoTime); err != nil {
		fmt.Println("视频信息获取失败！")
		fmt.Print(err)
		c.JSON(http.StatusOK, FeedResponse{
			Response:  model.Response{StatusCode: 1, StatusMsg: err.Error()},
			VideoList: videoList,
			NextTime:  nextTime,
		})
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
