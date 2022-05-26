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
   从视频中提取封面图片保存函数
   ps: 需要额外安装ffmpeg http://ffmpeg.org/download.html
   参数字段说明: 1.videoPath 要提取图片的视频路径; 2.toSavePath 封面图保存地址
*/
func videoToImage(videoPath string, toSavePath string) {
	arg := []string{"-hide_banner"}
	arg = append(arg, "-i", videoPath)
	arg = append(arg, "-r", "1")
	arg = append(arg, "-f", "image2")
	arg = append(arg, "-frames:v", "1") // 截取一张
	arg = append(arg, "-q", "8")        // 设置图片压缩等级，越高压缩越大
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

// GetVideoById 根据视频ID获取视频完整信息
func GetVideoById(videoId int64) (error, Video) {
	var video Video
	var dbVideoInfo DbVideoInfo
	// 视频信息
	if err := Db.Find(&dbVideoInfo, "video_id = ?", videoId).Error; err != nil {
		fmt.Println("视频信息查询失败！")
		return err, video
	}
	// 作者信息
	var dbUserInfo DbUserInfo
	if err := Db.Find(&dbUserInfo).Where("user_id = ?", dbVideoInfo.UserId).Error; err != nil {
		fmt.Println("视频作者信息查询失败！")
		return err, video
	}
	// 获取作者关注数和被关注数
	var follow Follow
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
	isFavorite = true

	// 点赞数
	var favoriteCount int64
	favoriteCount = FavoriteCount(videoId)
	fmt.Printf("\n视频点赞数：%d", favoriteCount)
	// 评论数
	var commentCount int64
	commentCount = CommentCount(videoId)
	fmt.Printf("\n视频评论数：%d", commentCount)
	// 拼接返回结果
	var author User
	author = User{Id: dbVideoInfo.UserId, Name: dbUserInfo.UserName, FollowCount: followCount, FollowerCount: fansCount, IsFollow: isFollow}
	video = Video{Id: videoId, PlayUrl: dbVideoInfo.PlayUrl, CoverUrl: dbVideoInfo.CoverUrl, Title: dbVideoInfo.Title,
		FavoriteCount: favoriteCount, CommentCount: commentCount, IsFavorite: isFavorite, Author: author}
	// To do: 更多异常处理考虑
	return nil, video
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
	// 测试根据视频id查询视频信息
	var videoData Video
	err, videoData := GetVideoById(1)
	if err == nil {
		fmt.Println("第1条视频查询结果：")
		fmt.Print(videoData)
	} else {
		// 异常处理
		fmt.Println("视频查询失败！")
		fmt.Print(err)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
