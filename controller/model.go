package controller

import "time"

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// DbUserInfo defines the structure that user informatiom is stored in database
type DbUserInfo struct {
	Id           int64
	UserName     string
	PasswordHash string
	Token        string
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

/*
type Follower struct {
	UserId     int64
	FollowerId int64
	Status     int64 //0 取关 1 关注 2 互关
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (DbUserInfo) TableName() string {
	return "UserName"
}

func (Follow) TableName() string {
	return "follow"
}
*/
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}