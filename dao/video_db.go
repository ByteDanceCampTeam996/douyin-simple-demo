package dao

import (
	"fmt"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"time"
)

// CreateNewVideo  新增视频接口
func CreateNewVideo(newVideo DbVideoInfo) error {
	err := Db.Create(newVideo).Error
	return err
}

// FindVideoInfoById 根据视频ID查询视频信息接口
func FindVideoInfoById(vid int64) (error, DbVideoInfo) {
	var videoInfo DbVideoInfo
	result := Db.Where("video_id = ?", vid).First(&videoInfo)
	return result.Error, videoInfo

}

// GetVideoList 首页获取视频列表，最多返回30条数据
func GetVideoList(token string, latestTime int64) (error, []Video, int64) {
	var activeUserId int64
	var nextTime int64
	// token解析
	if token == "" {
		fmt.Println("当前用户未登陆")
		activeUserId = -1 // 未登陆用户
	} else {
		if _, exist := UserExistByToken(token); exist {
			var activeUser DbUserInfo
			_, activeUser = FindUserByToken(token)
			activeUserId = activeUser.Id // 获取当前登陆用户id
		} else {
			activeUserId = -1
		}
	}
	// 将时间处理为时间戳
	formatTimeStr := time.Unix(latestTime/1000-60*60*24*7, 0).Format("2006-01-02 15:04:05.0001")

	fmt.Println(formatTimeStr)
	// 视频数据数组
	var videoList []Video
	var dbVideoTotal []DbVideoInfo
	// 限制了最多返回30条视频数据
	rows, err := Db.Model(&dbVideoTotal).Where("created_time > ?", formatTimeStr).Limit(30).Order("created_time desc").Rows()
	defer rows.Close()
	if err != nil {
		fmt.Println("视频数据查询失败")
		return err, videoList, time.Now().Unix()
	}
	// 按行查找
	for rows.Next() {
		var dbVideoInfo DbVideoInfo
		// ScanRows 将一行扫描
		Db.ScanRows(rows, &dbVideoInfo)
		// 获取单个video的video_id和作者user_id
		videoId := dbVideoInfo.VideoId
		userId := dbVideoInfo.UserId
		// 获取当前视频时间作为下一次传的参数
		nextTime = dbVideoInfo.CreatedTime.Unix()
		// 当前用户是否关注了作者，登陆与非登陆情况考虑
		var isFollow bool // 默认为false
		if activeUserId != -1 {
			var afollow Follow
			if err = Db.Where("user_id = ? AND follow_id = ?", activeUserId, dbVideoInfo.UserId).Find(&afollow).Error; err != nil {
				fmt.Println("查不到视频作者和当前用户的关注关系！")
			}
			if afollow.Status == 1 {
				isFollow = true
			}
		}
		// 当前用户是否点赞了视频,登陆与非登陆状态考虑
		var isFavorite bool
		if activeUserId != -1 {
			var afavourite DbFavorite
			if err = Db.Where("uid = ? AND vid = ? AND status = ?", activeUserId, videoId, 1).Find(&afavourite).Error; err != nil {
				fmt.Println("查不到当前用户和当前视频的喜欢关系！")
			}
			if afavourite.Status == 1 {
				isFavorite = true
			}
		}

		// 查看作者信息
		var dbUserInfo DbUserInfo
		if err = Db.Find(&dbUserInfo).Where("user_id = ?", userId).Error; err != nil {
			fmt.Println("视频作者信息查询失败！")
			return err, videoList, time.Now().Unix()
		}
		// 获取作者关注数和被关注数
		var follow Follow
		var followCount int64
		var fansCount int64
		if err = Db.Where("follow_id = ? AND status <> 1", dbVideoInfo.UserId).Find(&follow).Count(&fansCount).Error; err != nil {
			fmt.Println("查不到视频作者粉丝数！")
		}
		if err = Db.Where("user_id = ? AND status <> 1", dbVideoInfo.UserId).Find(&follow).Count(&followCount).Error; err != nil {
			fmt.Println("查不到视频作者关注数！")
		}
		// 点赞数
		var favoriteCount int64
		var dbFavorite DbFavorite
		if err = Db.Find(&dbFavorite).Where("vid = ? AND status = ?", videoId, 1).Count(&favoriteCount).Error; err != nil {
			fmt.Println("查不到视频点赞数！")
		}
		// 评论数
		var commentCount int64
		var dbComment DbComment
		if err = Db.Find(&dbComment).Where("vid = ?", videoId).Count(&commentCount).Error; err != nil {
			fmt.Println("查不到视频评论数！")
		}
		// 拼接返回结果
		var author User
		author = User{Id: dbVideoInfo.UserId, Name: dbUserInfo.UserName, FollowCount: followCount, FollowerCount: fansCount, IsFollow: isFollow}
		var video Video
		// 注意为空（false、0等）的时候某些字段不显示，存在非必需项
		video = Video{Id: videoId, PlayUrl: dbVideoInfo.PlayUrl, CoverUrl: dbVideoInfo.CoverUrl, Title: dbVideoInfo.Title,
			FavoriteCount: favoriteCount, CommentCount: commentCount, IsFavorite: isFavorite, Author: author}
		videoList = append(videoList, video)
	}
	// To do: 更多异常处理考虑
	return nil, videoList, nextTime
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
