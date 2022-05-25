package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DbComment struct {
	Id         int64 `gorm:"primary_key"`
	Vid        int64
	Content    string
	CreateDate string
	Uid        int64
	UserInfo   DbUserInfo `gorm:"ForeignKey:Uid"`
}
type CommentInfo struct {
	Id int64 `gorm:"column:id;`

	Content    string `gorm:"column:content"`
	CreateDate string `gorm:"column:create_date"`
	Uid        int64  `gorm:"column:uid"`
	UserName   string `gorm:"column:user_name"`
}

//var comments []comment

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentInsert(ct DbComment) {

	Db.Create(&ct)

}
func CommentFindByVid(vid int64) (comments []CommentInfo, err error) {
	//var com []Comment
	//Db.Debug().Where("vid=?", vid).Joins("db_user_infos").Find(&comments)
	//Db.Debug().Joins("db_user_infos").Find(&comments)
	//var results []map[string]interface{}
	Db.Debug().Model(&DbComment{}).Select("db_comments.id ,db_comments.content ,db_comments.create_date,db_comments.uid,db_user_infos.user_name").Joins("left join db_user_infos on db_user_infos.id = db_comments.uid").Scan(&comments)

	// defer rows.Close()
	// for rows.Next() {

	// 	var com CommentInfo
	// 	println("*****************************************")
	// 	println(rows.NextResultSet())
	// 	err = Db.ScanRows(rows, &com)
	// 	if err != nil {
	// 		println("序列化失败")
	// 		println(err)
	// 	} else {
	// 		//comments = append(comments, com)
	// 	}
	// 	return nil, err
	// }
	// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

	//res := Db.Where("vid=?", vid).Find(&comments)

	return
}

func CommentDeleteById(id int64) (err error) {
	res := Db.Where("id=?", id).Delete(DbComment{})
	err = res.Error
	return
}
func CommentAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := UserExistByToken(token); exist {
		action_type, err := strconv.Atoi(c.Query("action_type"))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type err"})
			return
		}
		_, dbUserInfo := FindUserByToken(token)
		uid := dbUserInfo.Id
		vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
		if err2 != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id format err"})
			return
		}

		if action_type == 1 {
			var cmt DbComment
			cmt.Content = c.Query("comment_text")
			cmt.CreateDate = time.Now().Format("12-05")
			cmt.Uid = uid
			cmt.Vid = vid
			CommentInsert(cmt)
		} else if action_type == 2 {
			id, err3 := strconv.ParseInt(c.Query("video_id"), 10, 64)
			if err3 != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id format err"})
				return
			}
			CommentDeleteById(id)
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	var com []Comment
	vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err2 != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id format err"})
		return
	}

	comments, _ := CommentFindByVid(vid)
	for i := 0; i < len(comments); i++ {
		user := User{Id: comments[i].Uid, Name: comments[i].UserName} //需要提供接口 ： getUserById(id int64)(User,error)
		com = append(com, Comment{
			comments[i].Id,
			user,
			comments[i].Content,
			comments[i].CreateDate,
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: com,
	})
}
