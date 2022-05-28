package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

// GetPublishVideoList 获取用户发布的视频列表
func GetPublishVideoList(token string, userId int64) (error, []Video) {
	// 要返回的视频数据数组
	var videoList []Video
	// 通过token判断用户是否登陆
	if _, exist := UserExistByToken(token); !exist {
		// 未登陆用户
		return nil, videoList
	}
	// 解析token
	_, dbUserInfo := FindUserByToken(token)
	if dbUserInfo.Id != userId {
		// 非法登陆，token和userId对应不上
		return nil, videoList
	}
	// 视频数据数组
	var dbVideoTotal []DbVideoInfo

	// 查看作者信息，解析token的时候已经查过一次作者信息了故直接复用
	// 当前用户是否关注了作者，自己发布的视频自己就是作者不会关注自己
	var isFollow bool // 默认为false
	// 获取作者关注数和被关注数
	var follow Follow
	var followCount int64
	var fansCount int64
	if err := Db.Where("follow_id = ? AND status <> 1", userId).Find(&follow).Count(&fansCount).Error; err != nil {
		fmt.Println("查不到视频作者粉丝数！")
	}
	if err := Db.Where("user_id = ? AND status <> 1", userId).Find(&follow).Count(&followCount).Error; err != nil {
		fmt.Println("查不到视频作者关注数！")
	}
	// 视频作者信息，自己发布的视频故作者信息只需查一次
	var author User
	author = User{Id: userId, Name: dbUserInfo.UserName, FollowCount: followCount, FollowerCount: fansCount, IsFollow: isFollow}

	// 查看用户发布的视频信息列表
	rows, err := Db.Model(&dbVideoTotal).Where("user_id = ?", userId).Order("created_time desc").Rows()
	defer rows.Close()
	if err != nil {
		fmt.Println("发布视频数据信息查询失败")
		return err, videoList
	}
	// 按行查找
	for rows.Next() {
		var dbVideoInfo DbVideoInfo
		// ScanRows 将一行扫描
		Db.ScanRows(rows, &dbVideoInfo)
		// 获取单个video的video_id和作者user_id
		videoId := dbVideoInfo.VideoId
		// 点赞数
		var favoriteCount int64
		var dbFavorite DbFavorite
		if err = Db.Find(&dbFavorite).Where("vid = ? AND status <> 1", videoId).Count(&favoriteCount).Error; err != nil {
			fmt.Println("查不到视频点赞数！")
		}
		// 评论数
		var commentCount int64
		var dbComment DbComment
		if err = Db.Find(&dbComment).Where("vid = ?", videoId).Count(&commentCount).Error; err != nil {
			fmt.Println("查不到视频评论数！")
		}
		// 是否点赞了视频
		var isFavorite bool
		if err = Db.Find(&dbFavorite).Where("uid = ? AND vid = ?", userId, videoId).Error; err != nil {
			fmt.Println("查不到当前用户和当前视频的点赞关系！")
		} else if dbFavorite.Status == 1 {
			isFavorite = true
		}
		// 拼接返回结果
		var video Video
		// 注意为空（false、0等）的时候某些字段不显示，存在非必需项
		video = Video{Id: videoId, PlayUrl: dbVideoInfo.PlayUrl, CoverUrl: dbVideoInfo.CoverUrl, Title: dbVideoInfo.Title,
			FavoriteCount: favoriteCount, CommentCount: commentCount, IsFavorite: isFavorite, Author: author}
		videoList = append(videoList, video)
	}
	// To do: 更多异常处理情况考虑
	return nil, videoList
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
		var bt3 bytes.Buffer
		bt3.WriteString(strings.Split(finalName, ".")[0])
		bt3.WriteString(".jpg")
		imgName := bt3.String()
		fmt.Printf("封面图名称:%s", imgName)
		imgSavePath = filepath.Join("./public", imgName)
		videoToImage(saveFile, imgSavePath)
		fmt.Printf("提取的封面图保存地址:%s", imgSavePath)
	}

	videoTitle := c.PostForm("title")

	// 修改为直接上传本地
	var saveVideoPath string
	var saveImgPath string
	// 分割获取要上传的视频名
	videoSlice := strings.Split(saveFile, "\\")
	videoName := videoSlice[len(videoSlice)-1]
	// 分割要上传上的封面图名
	imgSlice := strings.Split(imgSavePath, "\\")
	imgName := imgSlice[len(imgSlice)-1]
	// 返回视频和图片的访问url地址
	var bt1 bytes.Buffer
	bt1.WriteString("http://")
	bt1.WriteString(ipAddress)
	bt1.WriteString(":8080/static/")
	bt1.WriteString(videoName)
	var bt2 bytes.Buffer
	bt2.WriteString("http://")
	bt2.WriteString(ipAddress)
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
