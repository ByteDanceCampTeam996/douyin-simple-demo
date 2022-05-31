package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"sync/atomic"
	"time"

	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
)

// DbHashSalt 使用SHA-256算法将密码加盐之后哈希存储
func DbHashSalt(password string, salt string) string {
	hash1 := sha256.New()
	hash1.Write([]byte(password + salt))
	sum := hash1.Sum(nil)
	return hex.EncodeToString(sum)
}

// GetRandString 返回一个长度为10的随机字符串
func GetRandString() string {
	b := make([]byte, 10)
	rand.Read(b)
	return hex.EncodeToString(b)
}

//通过过滤一些基本的sql关键字防止sql注入
func IsLegalUserName(username string) bool {
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return !re.MatchString(username)
}

// UserAppend 通过用户名和密码添加一个新用户，返回用户的id和token
func UserAppend(username string, password string) (int64, string) {
	atomic.AddInt64(&UserIdSequence, 1)
	newUser := DbUserInfo{
		Id:           UserIdSequence,
		UserName:     username,
		PasswordHash: DbHashSalt(password, username),
		Token:        GetRandString(),
	}
	newFollowInfo := UserFollowInfo{
		UserId:        UserIdSequence,
		Name:          username,
		FollowCount:   0,
		FollowerCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	CreateNewUser(newUser)
	CreateNewFollowInfo(newFollowInfo)
	return UserIdSequence, newUser.Token
}
