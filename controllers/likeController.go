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

type LikeController interface {
	Add(ctx *gin.Context)
	ListUser(ctx *gin.Context)
}

type likeController struct {
	likeService services.LikeService
	jwtService  services.JWTService
}

func NewLikeController(likeServ services.LikeService, jwtServ services.JWTService) LikeController {
	return &likeController{
		likeService: likeServ,
		jwtService:  jwtServ,
	}
}

func (c *likeController) Add(ctx *gin.Context) {
	var reqLike *dto.RequestLike
	err := ctx.ShouldBind(&reqLike)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed get object like", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusUnprocessableEntity, res)
		return
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		reqLike.UserID = userID
		result, err := c.likeService.Add(ctx, reqLike)
		if err != nil {
			res := helpers.BuildErrorResponse("Failed to created like", err.Error(), helpers.EmptyObj{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		res := helpers.BuildResponse(true, "Success", result)
		ctx.JSON(http.StatusCreated, res)
	}
}

func (c *likeController) ListUser(ctx *gin.Context) {
	postID := ctx.Param("post_id")
	read, err := c.likeService.ListUser(ctx, postID)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed to read like user", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	} else {
		response := helpers.BuildResponse(true, "Readed!", read)
		ctx.JSON(http.StatusOK, response)
	}
}

func (c *likeController) getUserIDByToken(token string) string {
	Token, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := Token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
