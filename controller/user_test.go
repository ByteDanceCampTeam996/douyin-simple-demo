package controller

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
	url = "http://127.0.0.1:8080/douyin/user?token=1cc4087b02517bc6bc31"
	send_request(url, 2)
	url = "http://127.0.0.1:8080/douyin/user/register?username=test2&password=123456"
	send_request(url, 1)
	url = "http://127.0.0.1:8080/douyin/user/login?username=test2&password=123456"
	send_request(url, 1)
	url = "http://127.0.0.1:8080/douyin/user?token=cfb147de2be47ca879bd"
	send_request(url, 2)
}
