package model

import "time"

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// DbUserInfo 定义了用户的登录信息在数据库中的存储结构
type DbUserInfo struct {
	Id           int64
	UserName     string
	PasswordHash string
	Token        string
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// DbVideoInfo 定义了视频信息在数据库中的存储结构
type DbVideoInfo struct {
	VideoId     int64 `gorm:"column:video_id; primaryKey; not null; autoIncrement;"`
	UserId      int64 `gorm:"column:user_id; not null;"`
	PlayUrl     string
	CoverUrl    string
	Title       string
	CreatedTime time.Time
}

type UserFollowInfo struct {
	UserId        int64
	Name          string
	FollowCount   int64
	FollowerCount int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
type DbRelationAction struct {
	UserId     int64
	Token      string
	ToUserId   int64
	ActionType int64
}
type Follow struct {
	UserId    int64
	FollowId  int64
	Status    int64 //0 取关 1 关注 2 互关
	CreatedAt time.Time
	UpdatedAt time.Time
}
type DbComment struct {
	Id         int64 `gorm:"primary_key;"`
	Vid        int64
	Content    string
	CreateDate string
	Uid        int64
	//UserInfo   DbUserInfo `gorm:"ForeignKey:Uid"`
}

type DbFavorite struct {
	Uid    int64
	Vid    int64
	Status int
}

/*
func (DbUserInfo) TableName() string {
	return "UserName"
}

func (Follow) TableName() string {
	return "follow_relation"
}
func (UserFollowInfo) TableName() string {
	return "user_follow_info"
}
*/
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
