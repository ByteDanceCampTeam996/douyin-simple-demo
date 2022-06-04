package dao

import (
	"fmt"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"time"
)

// CreateNewVideo  新增视频接口
func CreateNewVideo(newVideo DbVideoInfo) error {
	return Db.Create(&newVideo).Error
}

// GetUserVideoList 获取用户已发布的全部视频
func GetUserVideoList(userId int64) (error, []DbVideoInfo) {
	var dbVideoTotal []DbVideoInfo
	if err := Db.Where("user_id = ?", userId).Order("created_time desc").Find(&dbVideoTotal).Error; err != nil {
		return err, dbVideoTotal
	}
	return nil, dbVideoTotal
}

// GetVideoList 获取首页视频列表
func GetVideoList(latestTime int64) (error, []DbVideoInfo) {
	// 将时间处理为时间戳
	formatTimeStr := time.Unix(latestTime/1000-60*60*24*7, 0).Format("2006-01-02 15:04:05.0001")
	fmt.Println(formatTimeStr)
	// 根据时间返回最多30条视频数据
	var dbVideoTotal []DbVideoInfo
	if err := Db.Where("created_time > ?", formatTimeStr).Limit(30).Order("created_time desc").Find(&dbVideoTotal).Error; err != nil {
		return err, dbVideoTotal
	}
	return nil, dbVideoTotal
}

// FindVideoInfoById 根据视频ID查询视频信息接口
func FindVideoInfoById(vid int64) (error, DbVideoInfo) {
	var videoInfo DbVideoInfo
	result := Db.Where("video_id = ?", vid).First(&videoInfo)
	// 没有找到记录时，会返回 ErrRecordNotFound 错误
	return result.Error, videoInfo
}
