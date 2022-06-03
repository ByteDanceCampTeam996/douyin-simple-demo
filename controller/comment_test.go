package controller

import (
	"testing"
)

func TestCommentAction(t *testing.T) {
	var url string
	url = "http://127.0.0.1:8080/douyin/comment/action/?token=5616b7cb5bf79ce32179&video_id=4&action_type=2&comment_id=29" //删除评论
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/comment/action/?token=5616b7cb5bf79ce32179&video_id=4&action_type=1&comment_text=你好" //添加评论
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/comment/action/?token=464fsa64sa4d32179&video_id=4&action_type=1&comment_text=你好" //错误token
	sendGetRequest(url)

}
func TestCommentList(t *testing.T) {

	var url string
	url = "http://127.0.0.1:8080/douyin/comment/list/?token=5616b7cb5bf79ce32179&video_id=4" //评论列表
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/comment/list/?token=5616b7cb5bf79ce32179&video_id=0" //请求不存在的评论列表
	sendGetRequest(url)

	url = "http://127.0.0.1:8080/douyin/favorite/list/?token=fsfsf5sffsfs5655&user_id=0" //错误token
	sendGetRequest(url)

}
