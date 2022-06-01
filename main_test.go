package main

import (
	"net/http/httptest"
	"testing"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Get(uri string, router *gin.Engine) *httptest.ResponseRecorder {
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}
func Post(uri string, router *gin.Engine) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应handler接口
	router.ServeHTTP(w, req)
	return w
}
func TestUser(t *testing.T) {
	r := gin.Default()

	initRouter(r)

	dao.ConnectDB()
	assert := assert.New(t)

	// 测试用户信息请求
	var url string
	url = "http://127.0.0.1:8080/douyin/user?token=1cc4087b02517bc6bc31"
	w := Get(url, r)
	assert.Equal(200, w.Code)
	// 测试用户注册请求
	url = "http://127.0.0.1:8080/douyin/user/register?username=test2&password=123456"
	w = Post(url, r)
	assert.Equal(200, w.Code)

}

/*
func send_request(url string, typ int) {
	if typ == 1 {
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
	} else {
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

}
//test feed
url = "http://127.0.0.1:8080/douyin/feed/"
send_request(url, 2)

//test user register and login
url = "http://127.0.0.1:8080/douyin/user?token=1cc4087b02517bc6bc31"
send_request(url, 2)
url = "http://127.0.0.1:8080/douyin/user/register?username=test2&password=123456"
send_request(url, 1)
url = "http://127.0.0.1:8080/douyin/user/login?username=test2&password=123456"
send_request(url, 1)
url = "http://127.0.0.1:8080/douyin/user?token=b379504dbde8c577bf59"
send_request(url, 2)
//test publish
url = "http://127.0.0.1:8080/douyin/publish/list/"
send_request(url, 2)
*/
