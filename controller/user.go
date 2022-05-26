package controller

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(0)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

/*
func FindUserByName(username string) DbUserInfo {
	var dbUserInfo DbUserInfo
	db.Where("user_name = ?", username).Find(&dbUserInfo)
	return dbUserInfo
}
*/

// DbHashSalt use sha-256 to hash the password with salt
func DbHashSalt(password string, salt string) string {
	hash1 := sha256.New()
	hash1.Write([]byte(password + salt))
	sum := hash1.Sum(nil)
	return hex.EncodeToString(sum)
}

// GetRandString returns a randomized string of fixed length
func GetRandString() string {
	b := make([]byte, 10)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func UserExistByToken(token string) (error, bool) {
	var dbUserInfo DbUserInfo
	if result := Db.Where("token = ?", token).First(&dbUserInfo); result.Error == nil {
		return nil, true
	} else {
		return result.Error, false
	}
}

func FindUserByToken(token string) (error, DbUserInfo) {
	var dbUserInfo DbUserInfo
	result := Db.Where("token = ?", token).Find(&dbUserInfo)
	if result.Error == nil {
		return nil, dbUserInfo
	} else {
		return result.Error, dbUserInfo
	}

}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var dbUserInfo DbUserInfo
	if result := Db.Where("user_name = ?", username).First(&dbUserInfo); result.Error == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := DbUserInfo{
			Id:           userIdSequence,
			UserName:     username,
			PasswordHash: DbHashSalt(password, username),
			Token:        GetRandString(),
		}
		newFollowInfo := UserFollowInfo{
			UserId:        userIdSequence,
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		Db.Create(newUser)
		Db.Create(newFollowInfo)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    newUser.Token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var dbUserInfo DbUserInfo
	if result := Db.Where("user_name = ?", username).First(&dbUserInfo); result.Error == nil {
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
		var userFollowInfo UserFollowInfo
		Db.Where("user_id = ?", dbUserInfo.Id).First(&userFollowInfo)
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
