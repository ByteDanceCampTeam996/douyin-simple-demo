package service

import (
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"sync"
)

func SetFavorite(uid int64, vid int64) error {
	return dao.FavoriteUpdate(uid, vid)
}

// GetFavoriteList 获取登陆用户的所有点赞视频列表
func GetFavoriteList(uid int64) (error, []model.Video) {
	// 要返回的视频列表
	var videoList []model.Video

	// 获取用户点赞的全部视频Id
	videoIds, err := dao.FavoriteVid(uid)
	if err != nil {
		return err, videoList
	}

	for i := 0; i < len(videoIds); i++ {
		err, dbVideoInfo := dao.FindVideoInfoById(videoIds[i])
		if err != nil {
			return err, videoList
		}
		videoId := dbVideoInfo.VideoId
		authorId := dbVideoInfo.UserId

		// 利用协程并行查询多表提高查询速度
		var wg sync.WaitGroup
		wg.Add(6)

		// 获取视频作者信息如用户名
		var authorName string
		go func() {
			defer wg.Done()
			_, userInfo := dao.FindUserInfoById(authorId)
			authorName = userInfo.UserName
		}()

		// 当前用户是否关注了作者
		var isFollow bool // 默认为false
		go func() {
			defer wg.Done()
			isFollow = dao.IsFollow(uid, authorId)
		}()

		// 当前用户是否点赞了视频,点赞列表默认为true
		var isFavorite bool
		isFavorite = true

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

		// 视频作者的关注数和粉丝数
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
		var author model.User
		author = model.User{Id: authorId, Name: authorName, FollowCount: followsCount, FollowerCount: fansCount, IsFollow: isFollow}
		var video model.Video
		video = model.Video{Id: videoId, PlayUrl: dbVideoInfo.PlayUrl, CoverUrl: dbVideoInfo.CoverUrl, Title: dbVideoInfo.Title,
			FavoriteCount: favoriteCount, CommentCount: commentCount, IsFavorite: isFavorite, Author: author}
		videoList = append(videoList, video)
	}

	return nil, videoList
}
