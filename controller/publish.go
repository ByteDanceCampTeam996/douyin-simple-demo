package controller

import (
	"bytes"
	"fmt"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	// 参数获取和校验
	token := c.PostForm("token")
	if _, exist := UserExistByToken(token); !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	_, dbUserInfo := FindUserByToken(token)
	t := time.Now()
	finalName := fmt.Sprintf("%d_%s_%s", dbUserInfo.Id, t.Format("20060102150405"), filename)
	saveFile := filepath.Join("./public", finalName)
	var imgSavePath string
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	} else {
		var bt3 bytes.Buffer
		bt3.WriteString(strings.Split(finalName, ".")[0])
		bt3.WriteString(".jpg")
		imgName := bt3.String()
		fmt.Printf("封面图名称:%s", imgName)
		imgSavePath = filepath.Join("./public", imgName)
		VideoToImage(saveFile, imgSavePath)
		fmt.Printf("提取的封面图保存地址:%s", imgSavePath)
	}

	videoTitle := c.PostForm("title")

	// 直接存储到本地，拼接返回的视频和图片url地址
	var saveVideoPath string
	var saveImgPath string
	//根据操作系统自动判断分隔符
	sysType := runtime.GOOS
	var sysSpliter string
	if sysType == "windows" {
		sysSpliter = "\\"
	} else {
		sysSpliter = "/"
	}
	// 分割获取要上传的视频名
	videoSlice := strings.Split(saveFile, sysSpliter)
	videoName := videoSlice[len(videoSlice)-1]
	// 分割获取要上传的封面图名
	imgSlice := strings.Split(imgSavePath, sysSpliter)
	imgName := imgSlice[len(imgSlice)-1]
	// 返回视频和图片的访问url地址
	var bt1 bytes.Buffer
	bt1.WriteString("http://")
	bt1.WriteString(IpAddress)
	bt1.WriteString(":8080/static/")
	bt1.WriteString(videoName)
	var bt2 bytes.Buffer
	bt2.WriteString("http://")
	bt2.WriteString(IpAddress)
	bt2.WriteString(":8080/static/")
	bt2.WriteString(imgName)
	// 获得拼接后的字符串
	saveVideoPath = bt1.String()
	saveImgPath = bt2.String()

	// 插入视频数据
	video := DbVideoInfo{UserId: dbUserInfo.Id, PlayUrl: saveVideoPath, CoverUrl: saveImgPath, Title: videoTitle, CreatedTime: time.Now()}
	if err := Db.Table("db_video_infos").Create(&video).Error; err == nil {
		fmt.Print("视频数据插入成功！")
	} else {
		fmt.Print("视频插入数据库失败！请检查sql语句！")
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
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
	var publishVideoList []Video
	if err, publishVideoList = GetPublishVideoList(token, newUserId); err != nil {
		fmt.Println("发布视频信息列表获取失败！")
		fmt.Print(err)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: publishVideoList,
	})
}
