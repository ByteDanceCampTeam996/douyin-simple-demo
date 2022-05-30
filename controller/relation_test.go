package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func sendPostRequest(url string) {
	contentType := "application/json"
	resp, err := http.Post(url, contentType, nil)
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Println(string(b))
}
func sendGetRequest(url string) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Println(string(b))
}
func TestRelationActionelation(t *testing.T) {
	/*
		关注操作 需要测试的：
		1、用户未登录
		关注
		2、A已关注B
		3、A取关B，B未关注A
		4、A取关B，B正在关注A
		5、A从未关注B，B未关注A
		6、A从未关注B，B正在关注A
		取关
		7、A未关注B
		8、B未关注A
		9、AB互关
	*/
	//根据数据库数据写？还是用mock？
	var url string
	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangsan&to_user_id=8&action_type=1"
	sendPostRequest(url)

	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangleidouyin&to_user_id=8&action_type=1"
	sendPostRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangleidouyin&to_user_id=3&action_type=1"
	sendPostRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangleidouyin&to_user_id=8&action_type=1"
	sendPostRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangleidouyin&to_user_id=2&action_type=1"
	sendPostRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangleidouyin&to_user_id=0&action_type=1"
	sendPostRequest(url)

	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangleidouyin&to_user_id=3&action_type=2"
	sendPostRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangleidouyin&to_user_id=2&action_type=2"
	sendPostRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/action/?token=zhangleidouyin&to_user_id=0&action_type=2"
	sendPostRequest(url)
}
func TestList(t *testing.T) {
	/*
		1、用户未登录
		2、无关注列表
		3、有关注列表
		4、数据库异常的情况？
	*/

	var url string
	url = "http://127.0.0.1:8080/douyin/relation/follow/list/?token=zhangsandouyin&user_id=2"
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/follow/list/?token=zhangleidouyin&user_id=1"
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/follow/list/?token=zhangleidouyin&user_id=0"
	sendGetRequest(url)

	url = "http://127.0.0.1:8080/douyin/relation/follower/list/?token=zhangsandouyin&user_id=2"
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/follower/list/?token=zhangleidouyin&user_id=0"
	sendGetRequest(url)
	url = "http://127.0.0.1:8080/douyin/relation/follower/list/?token=zhangleidouyin&user_id=1"
	sendGetRequest(url)

}
