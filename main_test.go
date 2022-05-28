package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

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
func TestUser(t *testing.T) {

	var url string
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
}
