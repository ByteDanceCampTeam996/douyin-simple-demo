package dao

import "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"

// CommentAction no practical effect, just check if token is valid
func CommentInsert(ct model.DbComment) {

	Db.Create(&ct)

}
func CommentFindByVid(vid int64) (comments []model.CommentInfo, err error) {

	Db.Debug().Model(&model.DbComment{}).Select("db_comments.id ,db_comments.content ,db_comments.create_date,db_comments.uid,db_user_infos.user_name").Joins("left join db_user_infos on db_user_infos.id = db_comments.uid").Where("vid=?", vid).Scan(&comments)

	return
}

func CommentDeleteById(id int64) (err error) {
	res := Db.Debug().Where("id=?", id).Delete(model.DbComment{})
	err = res.Error
	return
}
