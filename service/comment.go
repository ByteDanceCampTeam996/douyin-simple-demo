package service

import (
	"encoding/json"
	"time"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"github.com/garyburd/redigo/redis"
)

func GetCommentList(vid int64) (com []model.Comment) {
	var rds redis.Conn
	if dao.Rdspool != nil {

		rds = dao.Rdspool.Get() //从连接池，取一个链接
		defer rds.Close()       //函数运行结束 ，把连接放回连接池
		is_key_exit, _ := redis.Bool(rds.Do("EXISTS", vid))
		if is_key_exit {
			r, _ := redis.Bytes(rds.Do("get", vid))
			if len(r) > 0 {
				err := json.Unmarshal(r, &com)
				if err == nil {
					return
				}
			}

		}

	}
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
	if dao.Rdspool != nil {
		data, _ := json.Marshal(com)
		rds.Do("set", vid, data) //redis set命令
	}
	return

}
func delCommentKey(vid int64) {
	if dao.Rdspool != nil {
		rds := dao.Rdspool.Get() //从连接池，取一个链接
		defer rds.Close()        //函数运行结束 ，把连接放回连接池
		rds.Do("del", vid)       //redis set命令

	}
	return
}
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
		delCommentKey(vid)
		println("put:", vid)
	}
	dao.Db.Model(model.DbComment{}).First(&cmt)

	cm.User.Id = cmt.Id
	cm.User.Name = name

	cm.Id = cmt.Id

	cm.Content = cmt.Content
	cm.CreateDate = cmt.CreateDate

	return cm, err
}
func DelComment(id int64) (err error) {

	var cmt model.DbComment
	dao.Db.Model(model.DbComment{}).Where("id=?", id).First(&cmt)
	err = dao.CommentDeleteById(id)
	if err != nil {

		return
	} else {
		println("del:", cmt.Vid)
		delCommentKey(cmt.Vid)
	}

	return err
}
