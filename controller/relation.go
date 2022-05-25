package controller

import (
	"errors"
	"fmt"
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
	if user, exist := usersLoginInfo[token]; exist {
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
	if _, exist := usersLoginInfo[token]; exist {
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
	if _, exist := usersLoginInfo[token]; exist {
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
	fmt.Println("返回记录数：%d %d", rowf, rowe)

	if relationaction.ActionType == 1 {
		//关注操作 但是如果原先已经关注了 就提醒返回之前已关注 不再更新 （当然实际情况下应该不会这样的 只是现在加上这个判断
		if rowf == 1 && follow.Status != 0 {
			fmt.Println("已关注 不要再点了")
			return errors.New("已关注 请勿重复操作！")
		}
		//关注操作 在数据库粉丝列表查看这个对象是否关注user 如果关注了就要进行更新 为互关 没有关注就在粉丝列表加上这个数据
		//关注列表添加一条数据 并从上述获取是否是互关
		//如果之前关注过 然后取关了 就进行表更新 没有关注过就进行插入
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
			//rowf = 0之前没有关注过 那么就需要进行创建 然后需要查看一下对方是否关注过当前对象
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
		//取关操作 能进来说明user是关注了对方的
		/*
			follow表里找user 确保现在是关注状态 才能进行取关 不然返回报错 提醒已经取关了
			follow表里找取关对象 如果没找到就不用管；如果找到了，就更新为关注状态，
			follower表里找user 将状态更新为取关
			follower表里找取关对象 没找到就不用管；找到了状态更新为关注
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
	//所以还是用struct+select进行更新吧  然后就是指定更新表要用结构体不要用table，table也不会更新时间的！！！
	//更新单个值 设置modeljiuok；更新多个值 用model select！
	var tmpserch []Follow
	userfollows := Db.Where("user_id = ? AND status <> 0", relationaction.UserId).Find(&tmpserch).RowsAffected
	userfans := Db.Where("follow_id = ? AND status <> 0", relationaction.UserId).Find(&tmpserch).RowsAffected
	fmt.Println(relationaction.UserId, userfollows, userfans)
	Db.Table("users").Where("user_id = ?", relationaction.UserId).Select("follow_count", "follower_count").Updates(&UserFollowInfo{FollowCount: userfollows, FollowerCount: userfans})
	//更新map

	userfollows = Db.Where("user_id = ? AND status <> 0", relationaction.ToUserId).Find(&tmpserch).RowsAffected
	userfans = Db.Where("follow_id = ? AND status <> 0", relationaction.ToUserId).Find(&tmpserch).RowsAffected
	Db.Model(&UserFollowInfo{}).Where("user_id = ?", relationaction.ToUserId).Select("follow_count", "follower_count").Updates(&UserFollowInfo{FollowCount: userfollows, FollowerCount: userfans})
	//Db.Model(&model.UserInfo{}).Where("user_id = ?", relationaction.ToUserId).Updates(&model.UserInfo{FollowCount: userfollows, FollowerCount: userfans})
	//更新map

	fmt.Println(relationaction.ToUserId, userfollows, userfans)
	return nil
}

//功能设计的就是查看自己的列表
func GetFollowList(userId int64) ([]User, error) {
	var followlist []User
	var follow []Follow
	//要得到user的关注对象 存在数据并且是关注状态的
	res := Db.Where("user_id = ? AND status <> 0", userId).Find(&follow)
	row := res.RowsAffected
	var i int64
	var tmpUser User
	var tmpf User
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
	var follow []Follow
	//要得到user的粉丝 存在数据并且是关注状态的
	res := Db.Where("follow_id = ? AND status <> 0", userId).Find(&follow)
	row := res.RowsAffected
	var i int64
	var tmpUser User
	var tmpf User
	for i = 0; i < row; i++ {
		if follow[i].Status == 2 {
			tmpUser.IsFollow = true
		} else {
			tmpUser.IsFollow = false
		}
		tmpUser.Id = follow[i].UserId

		Db.Where("user_id = ?", tmpUser.Id).Find(&tmpf)
		tmpUser.FollowCount = tmpf.FollowCount
		tmpUser.FollowerCount = tmpf.FollowerCount
		tmpUser.Name = tmpf.Name
		followerlist = append(followerlist, tmpUser)
	}
	fmt.Println(followerlist)
	return followerlist, nil
}
