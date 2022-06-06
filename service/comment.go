package service

import (
	"encoding/json"
	"time"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"github.com/garyburd/redigo/redis"
)

//获取评论列表
func GetCommentList(vid int64) (com []model.Comment) {
	var rds redis.Conn
	if dao.Rdspool != nil { //判断redis连接池是否实例化

		rds = dao.Rdspool.Get() //从连接池，取一个链接
		defer rds.Close()       //函数运行结束 ，把连接放回连接池
		is_key_exit, _ := redis.Bool(rds.Do("EXISTS", vid))
		if is_key_exit { //如果对应视频id的key存在
			r, _ := redis.Bytes(rds.Do("get", vid))
			if len(r) > 0 {
				err := json.Unmarshal(r, &com) //反序列化
				if err == nil {
					return
				}
			}

		}

	}

	//根据视频id查询评论列表
	comments, _ := dao.CommentFindByVid(vid)

	for i := 0; i < len(comments); i++ {
		user := model.User{Id: comments[i].Uid, Name: comments[i].UserName}
		com = append(com, model.Comment{
			comments[i].Id,
			user,
			comments[i].Content,
			comments[i].CreateDate,
		})
	}
	//逻辑执行到这里，redis肯定是没有缓存的，所以这里缓存一下
	if dao.Rdspool != nil {
		data, _ := json.Marshal(com)
		rds.Do("set", vid, data)
	}
	return

}

//删除视频id映射的缓存
func delCommentKey(vid int64) {
	if dao.Rdspool != nil {
		rds := dao.Rdspool.Get() //从连接池，取一个链接
		defer rds.Close()        //函数运行结束 ，把连接放回连接池
		rds.Do("del", vid)

	}
	return
}

//添加评论
func PutComment(vid int64, content string, token string) (cm model.Comment, err error) {

	_, dbUserInfo := dao.FindUserByToken(token)
	uid := dbUserInfo.Id
	name := dbUserInfo.UserName

	var cmt model.DbComment
	cmt.Content = content
	cmt.CreateDate = time.Now().Format("01-02")
	cmt.Uid = uid
	cmt.Vid = vid
	res := dao.Db.Create(&cmt)
	if res.Error != nil {

		err = res.Error
		return
	} else {
		delCommentKey(vid) //新曾评论，redis里的缓存过期，将其删除

	}
	dao.Db.Model(model.DbComment{}).First(&cmt)

	cm.User.Id = cmt.Id
	cm.User.Name = name

	cm.Id = cmt.Id

	cm.Content = cmt.Content
	cm.CreateDate = cmt.CreateDate

	return cm, err
}

//删除评论
func DelComment(id int64) (err error) {

	var cmt model.DbComment
	dao.Db.Model(model.DbComment{}).Where("id=?", id).First(&cmt)
	err = dao.CommentDeleteById(id)
	if err != nil {

		return
	} else {

		delCommentKey(cmt.Vid) //redis里的缓存过期，将其删除
	}

	return err
}
