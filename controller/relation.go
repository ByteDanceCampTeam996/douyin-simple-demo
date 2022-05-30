package controller

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
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
	relationaction := DbRelationAction{
		Token:      token,
		ToUserId:   to_user_id,
		ActionType: action_type,
	}
	//用token判断当前用户是否存在？是否登录 不存在则提示用户登录或者注册 存在则进行关注|取关操作
	if _, exist := UserExistByToken(token); exist {
		_, user := FindUserByToken(token)
		relationaction.UserId = user.Id
		fmt.Println(relationaction)
		err := DoRelationAction(relationaction)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0, StatusMsg: "操作成功"})
		}

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	//获取query数据 然后调用service进行处理
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if _, exist := UserExistByToken(token); exist {
		userlist, err := GetFollowList(user_id)
		if err != nil {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  "获取成功",
				},
				//userlist是从数据库获取的数据
				UserList: userlist,
			})
		}

	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
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
		userlist, err := GetFollowerList(user_id)
		if err != nil {
			c.JSON(http.StatusNoContent, UserListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  "获取成功",
				},
				//userlist是从数据库获取的数据
				UserList: userlist,
			})
		}

	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "请先登录！",
			},
		})
	}
}

func DoRelationAction(relationaction DbRelationAction) error {
	fmt.Println(relationaction)
	if relationaction.UserId == relationaction.ToUserId {
		return errors.New("不能对自己进行操作！")
	}
	var follow Follow
	var follower Follow
	resf := Db.Where("user_id = ? AND follow_id = ?",
		relationaction.UserId, relationaction.ToUserId).Find(&follow)
	rese := Db.Where("user_id = ? AND follow_id = ?",
		relationaction.ToUserId, relationaction.UserId).Find(&follower)
	rowf := resf.RowsAffected
	rowe := rese.RowsAffected

	if (resf.Error != nil && errors.Is(resf.Error, gorm.ErrRecordNotFound) == false) || (rese.Error != nil && errors.Is(rese.Error, gorm.ErrRecordNotFound) == false) {
		return errors.New("Query error")
	}
	fmt.Printf("返回记录数：%d %d", rowf, rowe)

	if relationaction.ActionType == 1 {
		//关注操作 但是如果原先已经关注了 就提醒返回之前已关注 不再更新 （当然实际情况下应该不会这样的 只是现在加上这个判断
		if rowf == 1 && follow.Status != 0 {
			fmt.Println("已关注 不要再点了")
			return errors.New("已关注 请勿重复操作！")
		}
		//关注操作
		//数据库中存在当前用户是取关对象的状态；判断对方是否正在关注当前用户，关注则更新为互关，未关注则修改当前用户正在关注
		if rowf == 1 {
			if rowe == 0 || follower.Status == 0 {
				//对方没有关注过当前用户
				//或者对方取关了当前用户 那么只需要更新当前用户的信息即可 不用管被关注对象
				Db.Model(&Follow{}).Where("user_id = ? AND follow_id = ?",
					relationaction.UserId, relationaction.ToUserId).Update("status", 1)

			} else {
				//接下去就是对方关注了当前用户 那么就要更新为互关
				//关注列表 双方更新
				Db.Model(&Follow{}).Where("user_id = ? AND follow_id = ?",
					relationaction.UserId, relationaction.ToUserId).Update("status", 2)
				Db.Model(&Follow{}).Where("user_id = ? AND follow_id = ?",
					relationaction.ToUserId, relationaction.UserId).Update("status", 2)
			}
		} else {
			fmt.Println("原来无数据")
			//rowf = 0 当前用户之前没有关注过对方 那么就需要进行创建 然后需要查看一下对方是否关注过当前对象
			if rowe == 0 || follower.Status == 0 {
				//对方没有关注 或者取关了
				res := Db.Model(&Follow{}).Create(&Follow{UserId: relationaction.UserId, FollowId: relationaction.ToUserId, Status: 1})
				fmt.Println(res)
			} else {
				//对方是粉丝 更为互关
				Db.Model(&Follow{}).Create(&Follow{UserId: relationaction.UserId, FollowId: relationaction.ToUserId, Status: 2})
				Db.Model(&Follow{}).Where("user_id = ? AND follow_id = ?",
					relationaction.ToUserId, relationaction.UserId).Update("status", 2)

			}

		}
	} else {
		/*
			取关操作
			follow表里找user 确保现在是关注状态 才能进行取关 不然返回报错 提醒已经取关了
			follow表里找取关对象 如果没找到就不用管；如果找到了，就更新为关注状态
		*/
		if rowf == 0 || follow.Status == 0 {
			//user之前没有关注过当前对象 出bug了，操作不对
			fmt.Println("你没关注ta呀！")
			return errors.New("你没关注ta呀！")
		}
		//应该关注或者互关状态
		if rowe == 0 || follower.Status == 0 {
			//只是关注状态 对方没有关注user  只需要更新user的数据
			Db.Model(&Follow{}).Where("user_id = ? AND follow_id = ?",
				relationaction.UserId, relationaction.ToUserId).Update("status", 0)
			//return nil
		} else {
			//rowf == 1 rowe == 2 互关状态 user状态改为取关 对方改为关注
			Db.Model(&Follow{}).Where("user_id = ? AND follow_id = ?",
				relationaction.UserId, relationaction.ToUserId).Update("status", 0)
			Db.Model(&Follow{}).Where("user_id = ? AND follow_id = ?",
				relationaction.ToUserId, relationaction.UserId).Update("status", 1)
		}
	}
	fmt.Println("更新users表")
	//上面关注取关操作都处理完了 需要更新一下user和对象的信息 即更新users表里面的关注粉丝数量
	//要注意！！用struct更新时，默认不更新0值！！！要么用select选定字段，要么用map进行更新
	//但是用map进行更新的时候又不会更新时间？是为啥？
	//所以还是用struct+select进行更新吧  然后就是更新单列数据时指定更新表要用结构体不要用table，table也不会更新时间的！！！多列是可以的
	//更新单个值 设置model；更新多个值 用model select！
	var tmpserch []Follow
	userfollows := Db.Where("user_id = ? AND status <> 0", relationaction.UserId).Find(&tmpserch).RowsAffected
	userfans := Db.Where("follow_id = ? AND status <> 0", relationaction.UserId).Find(&tmpserch).RowsAffected
	fmt.Println(relationaction.UserId, userfollows, userfans)
	Db.Model(&UserFollowInfo{}).Where("user_id = ?", relationaction.UserId).Select("follow_count", "follower_count").Updates(&UserFollowInfo{FollowCount: userfollows, FollowerCount: userfans})


	userfollows = Db.Where("user_id = ? AND status <> 0", relationaction.ToUserId).Find(&tmpserch).RowsAffected
	userfans = Db.Where("follow_id = ? AND status <> 0", relationaction.ToUserId).Find(&tmpserch).RowsAffected
	Db.Model(&UserFollowInfo{}).Where("user_id = ?", relationaction.ToUserId).Select("follow_count", "follower_count").Updates(&UserFollowInfo{FollowCount: userfollows, FollowerCount: userfans})

	fmt.Println(relationaction.ToUserId, userfollows, userfans)
	return nil
}

//功能设计的就是查看自己的列表
func GetFollowList(userId int64) ([]User, error) {
	var followlist []User
	var follow []Follow
	//要得到user的关注对象 存在数据并且是关注状态的
	res := Db.Where("user_id = ? AND status <> 0", userId).Find(&follow)
	if res.Error != nil {
		return followlist, res.Error
	}
	row := res.RowsAffected
	var i int64
	//这个是用于返回的结构
	var tmpUser User
	//这个是存在数据库的用户关注粉丝数量信息
	var tmpf UserFollowInfo
	for i = 0; i < row; i++ {
		//因为是关注列表 那肯定是关注他了的
		tmpUser.IsFollow = true
		tmpUser.Id = follow[i].FollowId

		Db.Where("user_id = ?", tmpUser.Id).Find(&tmpf)
		tmpUser.FollowCount = tmpf.FollowCount
		tmpUser.FollowerCount = tmpf.FollowerCount
		tmpUser.Name = tmpf.Name
		followlist = append(followlist, tmpUser)
	}
	fmt.Println(followlist)
	return followlist, nil

}
func GetFollowerList(userId int64) ([]User, error) {
	//获取粉丝列表
	var followerlist []User
	var follower []Follow
	//要得到user的粉丝 存在数据并且是关注状态的
	res := Db.Where("follow_id = ? AND status <> 0", userId).Find(&follower)
	if res.Error != nil {
		return followerlist, res.Error
	}
	row := res.RowsAffected
	var i int64
	//这个是用于返回的结构
	var tmpUser User
	//这个是存在数据库的用户关注粉丝数量信息
	var tmpf UserFollowInfo
	for i = 0; i < row; i++ {
		//如果是互关状态 说明当前用户关注了这个粉丝；否则就是没关注
		if follower[i].Status == 2 {
			tmpUser.IsFollow = true
		} else {
			tmpUser.IsFollow = false
		}
		tmpUser.Id = follower[i].UserId

		Db.Where("user_id = ?", tmpUser.Id).Find(&tmpf)
		tmpUser.FollowCount = tmpf.FollowCount
		tmpUser.FollowerCount = tmpf.FollowerCount
		tmpUser.Name = tmpf.Name
		followerlist = append(followerlist, tmpUser)
	}
	fmt.Println(followerlist)
	return followerlist, nil
}
