package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

/*
	videoToImage 从视频中提取封面图片保存函数
    ps: 需要额外安装ffmpeg
    函数传的参数字段: 1.videoPath 要提取图片的视频路径; 2.toSavePath 封面图保存地址
*/
func videoToImage(videoPath string, toSavePath string) {
	arg := []string{"-hide_banner"}
	arg = append(arg, "-i", videoPath)
	arg = append(arg, "-r", "1")
	arg = append(arg, "-f", "image2")
	arg = append(arg, "-frames:v", "1")
	arg = append(arg, toSavePath)
	// 通过命令行运行ffmpeg截取视频帧图片保存为封面图
	cmd := exec.Command("ffmpeg", arg...)
	cmd.Stderr = os.Stderr
	fmt.Println("Run", cmd)
	err := cmd.Run()
	if err != nil {
		return
	}
	fmt.Println("提取视频封面图成功！")
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
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
		imgName := strings.Split(finalName, ".")[0] + ".jpg"
		fmt.Printf("封面图名称:%s", imgName)
		imgSavePath = filepath.Join("./public", imgName)
		go videoToImage(saveFile, imgSavePath)
		fmt.Printf("提取的封面图保存地址:%s", imgSavePath)
	}

	videoTitle := c.PostForm("title")

	// 插入视频数据
	video := DbVideoInfo{UserId: dbUserInfo.Id, PlayUrl: saveFile, CoverUrl: imgSavePath, Title: videoTitle, CreatedTime: time.Now()}
	if err := Db.Table("db_video_infos").Create(&video).Error; err == nil {
		fmt.Print("视频数据插入成功！")
	} else {
		fmt.Print("视频插入失败！请检查sql语句！")
		c.JSON(http.StatusOK, Response{
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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
