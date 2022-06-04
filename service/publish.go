package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func IsVideoType(filename string) bool {
	fileSlice := strings.Split(filename, ".")
	fileType := fileSlice[len(fileSlice)-1]
	// 根据后缀判断视频格式
	if fileType == "mp4" || fileType == "avi" || fileType == "mov" || fileType == "mpg" || fileType == "mpeg" || fileType == "wmv" || fileType == "rm" || fileType == "ram" {
		return true
	} else {
		return false
	}
}

// VideoToImage  从视频中提取封面图片并保存函数
// ps: 需要额外安装ffmpeg http://ffmpeg.org/download.html
func VideoToImage(videoPath string, toSavePath string) error {
	arg := []string{"-hide_banner"}
	arg = append(arg, "-i", videoPath)
	arg = append(arg, "-r", "1")
	arg = append(arg, "-f", "image2")
	arg = append(arg, "-frames:v", "1") // 截取一张
	arg = append(arg, "-q", "8")        // 设置图片压缩等级，越高压缩程度越大
	arg = append(arg, toSavePath)
	// 通过命令行运行ffmpeg截取视频帧图片保存为封面图
	cmd := exec.Command("ffmpeg", arg...)
	cmd.Stderr = os.Stderr
	fmt.Println("Run", cmd)
	err := cmd.Run()
	if err != nil {
		fmt.Println("提取视频封面图失败！")
		return err
	}
	fmt.Println("提取视频封面图成功！")
	return nil
}

// GetSavedUrlAddress 获取视频和图片的保存URL地址
func GetSavedUrlAddress(toSaveFilePath string) string {
	// 根据操作系统自动判断分隔符
	sysType := runtime.GOOS
	var sysSpliter string
	if sysType == "windows" {
		sysSpliter = "\\"
	} else {
		sysSpliter = "/"
	}
	// 分割获取文件名
	fileSlice := strings.Split(toSaveFilePath, sysSpliter)
	fileName := fileSlice[len(fileSlice)-1]
	// 拼接返回文件存储的访问URL地址
	var builder strings.Builder
	builder.WriteString("http://")
	builder.WriteString(IpAddress)
	builder.WriteString(":8080/static/")
	builder.WriteString(fileName)
	fileUrlAddress := builder.String()
	return fileUrlAddress
}

// VideoAppend 新增新视频到数据库
func VideoAppend(userId int64, playUrl string, coverUrl string, title string) error {
	newVideo := DbVideoInfo{UserId: userId, PlayUrl: playUrl, CoverUrl: coverUrl, Title: title, CreatedTime: time.Now()}
	return dao.CreateNewVideo(newVideo)
}

// GetUserPublishVideoList 获取登陆用户发布的全部视频列表
func GetUserPublishVideoList(token string, userId int64) (error, []Video) {
	// 要返回的视频数据数组
	var videoList []Video

	// 通过token是否存在判断用户是否登陆
	if _, exist := dao.UserExistByToken(token); !exist {
		// 未登陆用户
		return errors.New("token失效，用户未登陆！"), videoList
	}
	// 解析token
	_, dbUserInfo := dao.FindUserByToken(token)
	if dbUserInfo.Id != userId {
		// 非法登陆，token和当前登陆用户的userId对应不上
		return errors.New("非法登陆！"), videoList
	}

	authorId := dbUserInfo.Id

	// 利用协程并行查询多表提高查询速度
	var wg sync.WaitGroup
	wg.Add(3)

	// 获取用户发布的全部视频
	var dbVideoInfo []DbVideoInfo
	var videoErr error
	go func() {
		defer wg.Done()
		videoErr, dbVideoInfo = dao.GetUserVideoList(authorId)
	}()

	// 当前用户是否关注了作者，自己发布的视频自己就是作者不会关注自己
	var isFollow bool // 默认为false
	// 当前用户的关注数和粉丝数
	var followsCount int64
	var fansCount int64
	go func() {
		defer wg.Done()
		followsCount = dao.GetAuthorFollowsCount(authorId)
	}()
	go func() {
		defer wg.Done()
		fansCount = dao.GetAuthorFansCount(authorId)
	}()

	wg.Wait()
	if videoErr != nil {
		return videoErr, videoList
	}
	// 视频作者信息，当前用户就是作者只需查一次
	var author User
	author = User{Id: userId, Name: dbUserInfo.UserName, FollowCount: followsCount, FollowerCount: fansCount, IsFollow: isFollow}
	// 完善视频信息表
	for i := 0; i < len(dbVideoInfo); i++ {
		videoId := dbVideoInfo[i].VideoId

		// 利用协程并行查询多表提高查询速度
		var wg sync.WaitGroup
		wg.Add(3)

		var favoriteCount int64
		go func() {
			defer wg.Done()
			favoriteCount = dao.FavoriteCount(videoId)
		}()

		var commentCount int64
		go func() {
			defer wg.Done()
			commentCount = dao.CommentCount(videoId)
		}()

		var isFavorite bool
		go func() {
			defer wg.Done()
			isFavorite = dao.IsFavorite(authorId, videoId)
		}()

		wg.Wait()
		// 拼接视频列表返回结果
		var video Video
		video = Video{Id: videoId, PlayUrl: dbVideoInfo[i].PlayUrl, CoverUrl: dbVideoInfo[i].CoverUrl, Title: dbVideoInfo[i].Title,
			FavoriteCount: favoriteCount, CommentCount: commentCount, IsFavorite: isFavorite, Author: author}
		videoList = append(videoList, video)
	}
	return nil, videoList
}

// 上传文件到对象存储，在真机调试时存在上传大文件超时以及视频加载时间久的问题，故暂时还是使用本地存储的方式

// 七牛云存储密钥配置
var (
	accessKey  = "g532qClG4i74jUv5m8yFHhRK5uoLZMA4x9fHkgD1"
	secretKey  = "tOQBUdiYVhoPJN1ADeSKcS2XwvUZ_5MjNstVetPy"
	bucketName = "soplaying"
)

// UploadFileToObjectStore 上传视频和封面图文件到对象存储七牛云平台，输入要上传的文件地址，返回存储后的访问url地址
func UploadFileToObjectStore(toUploadVideoPath string, toUploadImgPath string) (error, string, string) {
	// 上传成功的时候返回存储后的视频和图片访问URL地址
	var videoUrlPath string
	var imgUrlPath string

	// 根据操作系统自动判断分隔符
	sysType := runtime.GOOS
	var sysSpliter string
	if sysType == "windows" {
		sysSpliter = "\\"
	} else {
		sysSpliter = "/"
	}
	// 分割获取要上传的视频名和图片名称
	videoSlice := strings.Split(toUploadVideoPath, sysSpliter)
	videoName := videoSlice[len(videoSlice)-1]
	imgSlice := strings.Split(toUploadImgPath, sysSpliter)
	imgName := imgSlice[len(imgSlice)-1]

	bucket := bucketName
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = true
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	// 上传视频
	if err := formUploader.PutFile(context.Background(), &ret, upToken, videoName, toUploadVideoPath, &putExtra); err != nil {
		return err, videoUrlPath, imgUrlPath
	}
	// 上传图片
	if err := formUploader.PutFile(context.Background(), &ret, upToken, imgName, toUploadImgPath, &putExtra); err != nil {
		return err, videoUrlPath, imgUrlPath
	}
	fmt.Println("文件上传成功！")
	fmt.Println(ret.Key, ret.Hash)

	// 拼接返回视频和图片的访问url地址
	var bt1 bytes.Buffer
	bt1.WriteString("http://rchn4rby8.hn-bkt.clouddn.com/") // 七牛云提供的临时解析域名即外链
	bt1.WriteString(videoName)
	videoUrlPath = bt1.String()
	var bt2 bytes.Buffer
	bt2.WriteString("http://rchn4rby8.hn-bkt.clouddn.com/")
	bt2.WriteString(imgName)
	imgUrlPath = bt2.String()
	return nil, videoUrlPath, imgUrlPath
}
