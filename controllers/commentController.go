package controllers

import (
	"fmt"
	"net/http"

	"github.com/aldisaputra17/posting-board/dto"
	"github.com/aldisaputra17/posting-board/helpers"
	"github.com/aldisaputra17/posting-board/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CommentController interface {
	Add(ctx *gin.Context)
	List(ctx *gin.Context)
	ListUserComment(ctx *gin.Context)
}

type commentController struct {
	commentService services.CommentService
	jwtService     services.JWTService
}

func NewCommentController(commentServ services.CommentService, jwtServ services.JWTService) CommentController {
	return &commentController{
		commentService: commentServ,
		jwtService:     jwtServ,
	}
}

func (c *commentController) Add(ctx *gin.Context) {
	var reqComment *dto.RequestComment
	err := ctx.ShouldBind(&reqComment)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed get object comment", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusUnprocessableEntity, res)
		return
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		reqComment.UserID = userID
		result, err := c.commentService.Add(ctx, reqComment)
		if err != nil {
			res := helpers.BuildErrorResponse("Failed to created comment", err.Error(), helpers.EmptyObj{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		res := helpers.BuildResponse(true, "Success", result)
		ctx.JSON(http.StatusCreated, res)
	}

}

func (c *commentController) List(ctx *gin.Context) {
	read, err := c.commentService.List(ctx)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed to read comment", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	} else {
		response := helpers.BuildResponse(true, "Readed!", read)
		ctx.JSON(http.StatusOK, response)
	}
}

func (c *commentController) ListUserComment(ctx *gin.Context) {
	postID := ctx.Param("post_id")
	read, err := c.commentService.ListUserComment(ctx, postID)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed to read like user", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	} else {
		response := helpers.BuildResponse(true, "Readed!", read)
		ctx.JSON(http.StatusOK, response)
	}
}

func (c *commentController) getUserIDByToken(token string) string {
	Token, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := Token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
