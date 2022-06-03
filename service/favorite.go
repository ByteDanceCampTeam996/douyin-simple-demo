package service

import (
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
)

func SetFavorite(uid int64, vid int64) error {
	return dao.FavoriteUpdate(uid, vid)
}

// FavoriteList all users have same favorite video list
func GetFavoriteList(uid int64) (videos []model.Video) {
	video_ids, err := dao.FavoriteVid(uid) // GetVideoById(videoId int64) Video
	if err == nil {
		println(video_ids)
		for i := 0; i < len(video_ids); i++ {
			_, video := dao.GetVideoById(uid, video_ids[i])
			videos = append(videos, video)
		}
	}

	return

}

//dao
