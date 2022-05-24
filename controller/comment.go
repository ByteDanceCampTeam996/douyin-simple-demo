package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type comment struct {
	Id         int64 `gorm:"-;primary_key;AUTO_INCREMENT"`
	Uid        int64
	Vid        int64
	Content    string
	CreateDate string
}

//var comments []comment

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentInsert(ct comment) {

	Db.Create(&ct)

}
func CommentFindByVid(vid int64) (comments []comment, err error) {
	res := Db.Where("vid=?", vid).Find(&comments)
	err = res.Error
	return
}
func CommentDeleteById(id int64) (err error) {
	res := Db.Where("id=?", id).Delete(comment{})
	err = res.Error
	return
}
func CommentAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		action_type, err := strconv.Atoi(c.Query("action_type"))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type err"})
			return
		}
		uid := usersLoginInfo[token].Id
		vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
		if err2 != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id format err"})
			return
		}
		// var id int64 = 1
		// if len(comments) > 0 {
		// 	id = comments[len(comments)-1].Uid + 1
		// }

		// if action_type == 1 {
		// 	comments = append(comments, comment{
		// 		id,
		// 		uid,
		// 		vid,
		// 		c.Query("comment_text"),
		// 		time.Now().Format("12-05"),
		// 	})
		// }
		if action_type == 1 {
			var cmt comment
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
	// for i := 0; i < len(comments); i++ {
	// 	if comments[i].Vid == vid {
	// 		com = append(com, Comment{
	// 			comments[i].Id,
	// 			usersLoginInfo[token],
	// 			comments[i].Content,
	// 			comments[i].CreateDate,
	// 		})
	// 	}
	// }
	comments, _ := CommentFindByVid(vid)
	for i := 0; i < len(comments); i++ {
		user := DemoUser //需要提供接口 ： getUserById(id int64)(User,error)
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
