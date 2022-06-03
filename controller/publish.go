package controller

import (
	"bytes"
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
	// To do: 视频标题或内容合法性校验

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	// 上传文件格式校验，判断是否为视频格式
	fileSlice := strings.Split(filename, ".")
	fileType := fileSlice[len(fileSlice)-1]
	if fileType != "mp4" && fileType != "avi" && fileType != "mov" && fileType != "mpg" && fileType != "mpeg" && fileType != "wmv" && fileType != "rm" && fileType != "ram" {
		c.JSON(http.StatusNotAcceptable, model.Response{
			StatusCode: 1,
			StatusMsg:  "上传视频格式不正确！",
		})
		return
	}

	// 添加时间和用户id重命名要保存的文件名称，避免文件名重复冲突
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

	// 重命名要保存的视频封面图名称
	var bt bytes.Buffer
	bt.WriteString(strings.Split(finalName, ".")[0])
	bt.WriteString(".jpg")
	imgName := bt.String()
	fmt.Printf("封面图名称:%s", imgName)
	imgSavePath := filepath.Join("./public", imgName)
	// 提取视频封面图
	service.VideoToImage(saveFile, imgSavePath)

	// 返回存储后的URL地址
	savedVideoPath := service.GetSavedUrlAddress(saveFile)
	savedImgPath := service.GetSavedUrlAddress(imgSavePath)

	// 插入视频数据
	newVideo := model.DbVideoInfo{UserId: dbUserInfo.Id, PlayUrl: savedVideoPath, CoverUrl: savedImgPath, Title: videoTitle, CreatedTime: time.Now()}
	if err := dao.CreateNewVideo(newVideo); err != nil {
		fmt.Print("视频插入数据库失败！请检查sql语句！")
		c.JSON(http.StatusInternalServerError, model.Response{
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
	userId := c.Query("user_id")
	token := c.Query("token")
	// string转成int64
	newUserId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Println("userId转换为int64失败，字段非法！")
	}
	var publishVideoList []model.Video
	if err, publishVideoList = dao.GetPublishVideoList(token, newUserId); err != nil {
		fmt.Println("发布视频信息列表获取失败！")
		fmt.Print(err)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: publishVideoList,
	})
}
