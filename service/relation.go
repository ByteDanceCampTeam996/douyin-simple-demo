package service

import (
	"errors"
	"fmt"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
)

func DoRelationAction(relationaction model.DbRelationAction) error {
	fmt.Println(relationaction)
	if relationaction.UserId == relationaction.ToUserId {
		return errors.New("不能对自己进行操作！")
	}
	var follow model.Follow
	var follower model.Follow
	follow, rowf := dao.GetFollowRecord(relationaction.UserId, relationaction.ToUserId)
	follower, rowe := dao.GetFollowRecord(relationaction.ToUserId, relationaction.UserId)
	fmt.Println(follow, follower)
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
				dao.UpdateFollowStatus(relationaction.UserId, relationaction.ToUserId, 1)
			} else {
				//接下去就是对方关注了当前用户 那么就要更新为互关
				//关注列表 双方更新
				dao.UpdateFollowStatus(relationaction.UserId, relationaction.ToUserId, 2)
				dao.UpdateFollowStatus(relationaction.ToUserId, relationaction.UserId, 2)

			}
		} else {
			fmt.Println("原来无数据")
			//rowf = 0 当前用户之前没有关注过对方 那么就需要进行创建 然后需要查看一下对方是否关注过当前对象
			if rowe == 0 || follower.Status == 0 {
				//对方没有关注 或者取关了
				dao.CreateFollow(relationaction.UserId, relationaction.ToUserId, 1)
			} else {
				//对方是粉丝 更为互关
				dao.CreateFollow(relationaction.UserId, relationaction.ToUserId, 2)
				dao.UpdateFollowStatus(relationaction.ToUserId, relationaction.UserId, 2)

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
			dao.UpdateFollowStatus(relationaction.UserId, relationaction.ToUserId, 0)
		} else {
			//rowf == 1 rowe == 2 互关状态 user状态改为取关 对方改为关注
			dao.UpdateFollowStatus(relationaction.UserId, relationaction.ToUserId, 0)
			dao.UpdateFollowStatus(relationaction.ToUserId, relationaction.UserId, 1)
		}
	}
	fmt.Println("更新users表")
	//上面关注取关操作都处理完了 需要更新一下user和对象的信息 即更新users表里面的关注粉丝数量
	//要注意！！用struct更新时，默认不更新0值！！！要么用select选定字段，要么用map进行更新
	//但是用map进行更新的时候又不会更新时间？是为啥？
	//所以还是用struct+select进行更新吧  然后就是更新单列数据时指定更新表要用结构体不要用table，table也不会更新时间的！！！多列是可以的
	//更新单个值 设置model；更新多个值 用model select！
	var tmpserch []model.Follow
	tmpserch, _ = dao.GetFollows(relationaction.UserId)
	userfollows := int64(len(tmpserch))
	tmpserch, _ = dao.GetFans(relationaction.UserId)
	userfans := int64(len(tmpserch))
	fmt.Println(relationaction.UserId, userfollows, userfans)
	dao.UpdateUserFollowInfo(relationaction.UserId, userfollows, userfans)

	tmpserch, _ = dao.GetFollows(relationaction.ToUserId)
	userfollows = int64(len(tmpserch))
	tmpserch, _ = dao.GetFans(relationaction.ToUserId)
	userfans = int64(len(tmpserch))
	dao.UpdateUserFollowInfo(relationaction.ToUserId, userfollows, userfans)

	fmt.Println(relationaction.ToUserId, userfollows, userfans)
	return nil
}

//功能设计的就是查看自己的列表
func GetFollowList(userId int64) ([]model.User, error) {
	var followlist []model.User
	var follow []model.Follow
	//要得到user的关注对象 存在数据并且是关注状态的
	//res := dao.Db.Where("user_id = ? AND status <> 0", userId).Find(&follow)
	follow, err := dao.GetFollows(userId)
	if err != nil {
		return followlist, err
	}
	row := len(follow)
	//这个是用于返回的结构
	var tmpUser model.User
	//这个是存在数据库的用户关注粉丝数量信息
	var tmpf model.UserFollowInfo
	for i := 0; i < row; i++ {
		//因为是关注列表 那肯定是关注他了的
		tmpUser.IsFollow = true
		tmpUser.Id = follow[i].FollowId

		//dao.Db.Where("user_id = ?", tmpUser.Id).Find(&tmpf)
		tmpf = dao.GetUserFollowInfo(tmpUser.Id)
		tmpUser.FollowCount = tmpf.FollowCount
		tmpUser.FollowerCount = tmpf.FollowerCount
		tmpUser.Name = tmpf.Name
		followlist = append(followlist, tmpUser)
	}
	fmt.Println(followlist)
	return followlist, nil

}
func GetFollowerList(userId int64) ([]model.User, error) {
	//获取粉丝列表
	var followerlist []model.User
	var follower []model.Follow
	//要得到user的粉丝 存在数据并且是关注状态的
	//res := dao.Db.Where("follow_id = ? AND status <> 0", userId).Find(&follower)
	follower, err := dao.GetFans(userId)
	if err != nil {
		return followerlist, err
	}
	row := len(follower)
	//这个是用于返回的结构
	var tmpUser model.User
	//这个是存在数据库的用户关注粉丝数量信息
	var tmpf model.UserFollowInfo
	for i := 0; i < row; i++ {
		//如果是互关状态 说明当前用户关注了这个粉丝；否则就是没关注
		if follower[i].Status == 2 {
			tmpUser.IsFollow = true
		} else {
			tmpUser.IsFollow = false
		}
		tmpUser.Id = follower[i].UserId

		//dao.Db.Where("user_id = ?", tmpUser.Id).Find(&tmpf)
		tmpf = dao.GetUserFollowInfo(tmpUser.Id)
		tmpUser.FollowCount = tmpf.FollowCount
		tmpUser.FollowerCount = tmpf.FollowerCount
		tmpUser.Name = tmpf.Name
		followerlist = append(followerlist, tmpUser)
	}
	fmt.Println(followerlist)
	return followerlist, nil
}
