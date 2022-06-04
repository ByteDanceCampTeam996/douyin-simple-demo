package controller

import (
	"net/http"
	"strconv"

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

		vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)

		if err2 != nil {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "video_id format err"})
			return
		}

		if action_type == 1 {

			Content := c.Query("comment_text")

			//判断评论是否合法，防止可能的sql注入
			if !service.IsLegalUserName(Content) {
				c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "Not a legal comment"})
				return
			}
			cm, err := service.PutComment(vid, Content, token)
			if err == nil {

				c.JSON(http.StatusOK, CommentResponse{
					Response:  model.Response{StatusCode: 0},
					MyComment: cm,
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
			service.DelComment(id)

		}

		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	vid, err2 := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err2 != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "video_id format err"})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: service.GetCommentList(vid),
	})
}
