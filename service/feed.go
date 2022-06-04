package service

import (
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"sync"
	"time"
)

// GetHomeVideoList 首页获取视频列表，token和latestTime均为可选项，最多返回30条数据
func GetHomeVideoList(token string, latestTime int64) (error, []Video, int64) {
	// 要返回的视频列表
	var videoList []Video

	var activeUserId int64 // 当前用户Id
	var nextTime int64     // 下一次请求的时间戳

	// token解析
	if token == "" {
		activeUserId = -1 // 未登陆用户标记
	} else {
		if _, exist := dao.UserExistByToken(token); exist {
			var activeUser DbUserInfo
			_, activeUser = dao.FindUserByToken(token)
			activeUserId = activeUser.Id // 获取当前登陆用户id
		} else {
			activeUserId = -1 // token失效，未登陆用户标记
		}
	}

	// 获取首页视频
	var dbVideoInfo []DbVideoInfo
	var videoErr error
	if videoErr, dbVideoInfo = dao.GetVideoList(latestTime); videoErr != nil {
		return videoErr, videoList, time.Now().Unix()
	}

	// 完善视频信息表
	for i := 0; i < len(dbVideoInfo); i++ {
		videoId := dbVideoInfo[i].VideoId
		authorId := dbVideoInfo[i].UserId

		// 利用协程并行查询多表提高查询速度
		var wg sync.WaitGroup
		wg.Add(7)

		// 获取视频作者信息如用户名
		var authorName string
		go func() {
			defer wg.Done()
			_, userInfo := dao.FindUserInfoById(authorId)
			authorName = userInfo.UserName
		}()

		// 当前用户是否关注了作者，登陆与非登陆情况考虑
		var isFollow bool // 默认为false
		go func() {
			defer wg.Done()
			if activeUserId != -1 {
				isFollow = dao.IsFollow(activeUserId, authorId)
			}
		}()

		// 当前用户是否点赞了视频,登陆与非登陆状态考虑
		var isFavorite bool
		go func() {
			defer wg.Done()
			if activeUserId != -1 {
				isFavorite = dao.IsFavorite(activeUserId, videoId)
			}
		}()

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

		// 视频作者的粉丝数和关注数
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
		// 拼接视频列表返回结果
		var author User
		author = User{Id: authorId, Name: authorName, FollowCount: followsCount, FollowerCount: fansCount, IsFollow: isFollow}
		var video Video
		video = Video{Id: videoId, PlayUrl: dbVideoInfo[i].PlayUrl, CoverUrl: dbVideoInfo[i].CoverUrl, Title: dbVideoInfo[i].Title,
			FavoriteCount: favoriteCount, CommentCount: commentCount, IsFavorite: isFavorite, Author: author}
		videoList = append(videoList, video)

		// 视频倒序排序，获取当前视频时间作为下一次传的参数
		nextTime = dbVideoInfo[i].CreatedTime.Unix()
	}
	return nil, videoList, nextTime

}
