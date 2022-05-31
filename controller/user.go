package controller

import (
	"net/http"

	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/service"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//判断用户名是否合法，防止可能的sql注入
	if !IsLegalUserName(username) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Not a legal username"},
		})
	}
	//判断用户名是否已经存在
	if _, exist := UserExistByUsername(username); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		//如果用户名不存在，添加新用户
		userid, token := UserAppend(username, password)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userid,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//判断用户名是否合法，防止可能的sql注入
	if !IsLegalUserName(username) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Not a legal username"},
		})
	}

	//判断用户名是否已经存在
	if _, exist := UserExistByUsername(username); exist {
		//重新计算哈希以验证密码是否正确
		_, dbUserInfo := FindUserByUsername(username)
		if DbHashSalt(password, username) == dbUserInfo.PasswordHash {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   dbUserInfo.Id,
				Token:    dbUserInfo.Token,
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Wrong Password"},
			})
		}
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {

	token := c.Query("token")

	if _, exist := UserExistByToken(token); exist {
		_, dbUserInfo := FindUserByToken(token)
		_, userFollowInfo := FindFollowInfoById(dbUserInfo.Id)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: User{
				Id:            dbUserInfo.Id,
				Name:          dbUserInfo.UserName,
				FollowCount:   userFollowInfo.FollowCount,
				FollowerCount: userFollowInfo.FollowerCount,
			},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
