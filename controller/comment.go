package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/service"
	"github.com/gin-gonic/gin"
)

//var comments []comment

type CommentListResponse struct {
	model.Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}
type CommentResponse struct {
	model.Response
	MyComment model.Comment `json:"comment,omitempty"`
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := dao.UserExistByToken(token); exist {
		action_type, err := strconv.Atoi(c.Query("action_type"))
		if err != nil {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "action_type err"})
			return
		}
		_, dbUserInfo := dao.FindUserByToken(token)
		uid := dbUserInfo.Id
		vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
		if err2 != nil {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "video_id format err"})
			return
		}

		if action_type == 1 {
			var cmt model.DbComment
			cmt.Content = c.Query("comment_text")

			//判断评论是否合法，防止可能的sql注入
			if !service.IsLegalUserName(cmt.Content) {
				c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "Not a legal comment"})
				return
			}
			cmt.CreateDate = time.Now().Format("01-02")
			cmt.Uid = uid
			cmt.Vid = vid
			res := dao.Db.Create(&cmt)

			var user model.User
			user.Id = cmt.Id
			user.Name = dbUserInfo.UserName
			if res.Error == nil {
				dao.Db.Model(model.DbComment{}).First(&cmt)
				println(cmt.Id)
				c.JSON(http.StatusOK, CommentResponse{
					Response: model.Response{StatusCode: 0},
					MyComment: model.Comment{
						cmt.Id,
						user,
						cmt.Content,
						cmt.CreateDate,
					},
				})
				return

			}

			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "comment fail"})
		} else if action_type == 2 {
			id, err3 := strconv.ParseInt(c.Query("comment_id"), 10, 64)
			if err3 != nil {
				c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "video_id format err"})
				return
			}
			dao.CommentDeleteById(id)
		}

		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	var com []model.Comment
	vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err2 != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "video_id format err"})
		return
	}

	comments, _ := dao.CommentFindByVid(vid)
	for i := 0; i < len(comments); i++ {
		user := model.User{Id: comments[i].Uid, Name: comments[i].UserName} //需要提供接口 ： getUserById(id int64)(User,error)
		com = append(com, model.Comment{
			comments[i].Id,
			user,
			comments[i].Content,
			comments[i].CreateDate,
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: com,
	})
}
