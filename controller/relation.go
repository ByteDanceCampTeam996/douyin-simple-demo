package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/service"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	model.Response
	UserList []model.User `json:"user_list"`
}

/*
数据存储crud处理在model里面
service，处理业务 例如调用数据库之类的
controller处理和外部数据的交互，简单来说就是获取url的数据，检查对不对，然后调用service进行下一步的处理
*/
/*
关系服务的话三个功能
关注或者取关用户
首先要得知被操作的对象id，然后从数据库中去寻找当前用户是否关注了当前用户
如果是在关注列表里面，就要进行取关操作：从关注列表移除
如果不在关注列表，就进行关注操作，加入到关注列表
获取关注列表
获取粉丝列表
*/
//现在只用以一个关注表来实现
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	//获取query数据 然后调用service进行处理  user_id根据token获取
	//user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	to_user_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	relationaction := model.DbRelationAction{
		Token:      token,
		ToUserId:   to_user_id,
		ActionType: action_type,
	}
	//用token判断当前用户是否存在？是否登录 不存在则提示用户登录或者注册 存在则进行关注|取关操作
	if _, exist := UserExistByToken(token); exist {
		_, user := FindUserByToken(token)
		relationaction.UserId = user.Id
		fmt.Println(relationaction)
		err := service.DoRelationAction(relationaction)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, model.Response{
				StatusCode: 0, StatusMsg: "操作成功"})
		}

	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	//获取query数据 然后调用service进行处理
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if _, exist := UserExistByToken(token); exist {
		userlist, err := service.GetFollowList(user_id)
		if err != nil {
			c.JSON(http.StatusOK, UserListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, UserListResponse{
				Response: model.Response{
					StatusCode: 0,
					StatusMsg:  "获取成功",
				},
				//userlist是从数据库获取的数据
				UserList: userlist,
			})
		}

	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "请先登录！",
			},
		})
	}
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	//获取query数据 然后调用service进行处理
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if _, exist := UserExistByToken(token); exist {
		userlist, err := service.GetFollowerList(user_id)
		if err != nil {
			c.JSON(http.StatusNoContent, UserListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, UserListResponse{
				Response: model.Response{
					StatusCode: 0,
					StatusMsg:  "获取成功",
				},
				//userlist是从数据库获取的数据
				UserList: userlist,
			})
		}

	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "请先登录！",
			},
		})
	}
}
