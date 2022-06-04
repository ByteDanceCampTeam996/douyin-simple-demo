package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"time"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/service"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	// 参数获取和校验
	token := c.PostForm("token")
	if _, exist := dao.UserExistByToken(token); !exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	videoTitle := c.PostForm("title")
	// To do: 视频标题或内容简介合法性校验，如长度限制、敏感词过滤等

	// 获取用户上传的文件数据
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 上传文件格式校验，判断是否为视频格式
	filename := filepath.Base(data.Filename)
	if !service.IsVideoType(filename) {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "上传视频格式不正确！",
		})
		return
	}

	// 添加用户id、当前时间重命名要保存的文件名称，避免文件名重复冲突
	_, dbUserInfo := dao.FindUserByToken(token)
	t := time.Now()
	finalName := fmt.Sprintf("%d_%s_%s", dbUserInfo.Id, t.Format("20060102150405"), filename)
	saveFile := filepath.Join("./public", finalName)
	// 将上传视频文件保存到本地public文件夹下
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 设置要保存的视频封面图名称及保存路径
	var bt strings.Builder
	bt.WriteString(strings.Split(finalName, ".")[0])
	bt.WriteString(".jpg")
	imgName := bt.String()
	imgSavePath := filepath.Join("./public", imgName)

	// 提取视频封面图，由于文件存储具有顺序性故串行化无使用并发。 To do: 后续可以考虑加入消息队列做异步操作
	if err := service.VideoToImage(saveFile, imgSavePath); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 返回存储后的URL地址
	savedVideoPath := service.GetSavedUrlAddress(saveFile)
	savedImgPath := service.GetSavedUrlAddress(imgSavePath)

	// 插入视频数据到数据库
	if err := service.VideoAppend(dbUserInfo.Id, savedVideoPath, savedImgPath, videoTitle); err != nil {
		fmt.Println("新增视频数据失败")
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// 参数获取和检验
	userId := c.Query("user_id")
	token := c.Query("token")

	// user_id格式处理： string转成int64
	newUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, model.Response{
			StatusCode: 1,
			StatusMsg:  "user_id字段格式非法！",
		})
		return
	}

	var publishVideoList []model.Video
	if err, publishVideoList = service.GetUserPublishVideoList(token, newUserId); err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			VideoList: publishVideoList,
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: publishVideoList,
	})
}
