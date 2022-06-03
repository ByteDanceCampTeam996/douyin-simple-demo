package controller

import (
	"testing"
)

func TestFavoriteAction(t *testing.T) {
	var url string
	url = "http://127.0.0.1:8080/douyin/favorite/action/?token=5616b7cb5bf79ce32179&video_id=4&action_type=1" //点赞
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/favorite/action/?token=5616b7cb5bf79ce32179&video_id=4&action_type=2" //取消点赞
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/favorite/action/?token=5616b7cb5bf79ce32179&video_id=0&action_type=1" //操作的视频不存在
	sendGetRequest(url)

}
func TestFavoriteList(t *testing.T) {

	var url string
	url = "http://127.0.0.1:8080/douyin/favorite/list/?token=5616b7cb5bf79ce32179&user_id=1" //正确请求
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/favorite/list/?token=5616b7cb5bf79ce32179&user_id=0" //不存在的user_id
	sendGetRequest(url)

	url = "http://127.0.0.1:8080/douyin/favorite/list/?token=fsfsf5sffsfs5655&user_id=0" //错误token
	sendGetRequest(url)

}
