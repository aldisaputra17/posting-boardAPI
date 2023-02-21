package controllers

import (
	"fmt"
	"net/http"

	"github.com/aldisaputra17/posting-board/dto"
	"github.com/aldisaputra17/posting-board/entities"
	"github.com/aldisaputra17/posting-board/helpers"
	"github.com/aldisaputra17/posting-board/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostController interface {
	Add(ctx *gin.Context)
	Deleted(ctx *gin.Context)
	List(ctx *gin.Context)
	ListPublic(ctx *gin.Context)
	CommentMost(ctx *gin.Context)
	LikeMost(ctx *gin.Context)
	FindByUserID(ctx *gin.Context)
}

type postController struct {
	postService services.PostService
	jwtService  services.JWTService
}

func NewPostController(postServ services.PostService, jwtServ services.JWTService) PostController {
	return &postController{
		postService: postServ,
		jwtService:  jwtServ,
	}
}

func (c *postController) Add(ctx *gin.Context) {
	var reqPost *dto.RequestPost
	err := ctx.ShouldBind(&reqPost)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed get object post", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusUnprocessableEntity, res)
		return
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		reqPost.UserID = userID
		result, err := c.postService.Add(ctx, reqPost)
		if err != nil {
			res := helpers.BuildErrorResponse("Failed to created post", err.Error(), helpers.EmptyObj{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		res := helpers.BuildResponse(true, "Success", result)
		ctx.JSON(http.StatusCreated, res)
	}
}

func (c *postController) Deleted(ctx *gin.Context) {
	var post entities.Post
	id := ctx.Param("id")
	post.ID = uuid.Must(uuid.Parse(id))
	authHeader := ctx.GetHeader("Authorization")
	token, errTkn := c.jwtService.ValidateToken(authHeader)
	if errTkn != nil {
		panic(errTkn)
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.postService.IsAllowedToEdit(ctx, id, userID) {
		err := c.postService.Deleted(ctx, post)
		if err != nil {
			res := helpers.BuildErrorResponse("Fail deleted post", err.Error(), helpers.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		}
		res := helpers.BuildResponse(true, "Ok", helpers.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
	} else {
		response := helpers.BuildErrorResponse("You dont have permission", "You are not the owner", helpers.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
	}
}

func (c *postController) List(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	list, err := c.postService.ListPrivate(ctx, userID)
	if err != nil {
		res := helpers.BuildErrorResponse("Fail fetch list post", err.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := helpers.BuildResponse(true, "Ok", list)
	ctx.JSON(http.StatusOK, res)
}

func (c *postController) ListPublic(ctx *gin.Context) {
	read, err := c.postService.ListPublic(ctx)
	if err != nil {
		response := helpers.BuildErrorResponse("Failed to readed", err.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	} else {
		response := helpers.BuildResponse(true, "Readed!", read)
		ctx.JSON(http.StatusOK, response)
	}
}

func (c *postController) CommentMost(ctx *gin.Context) {
	code := http.StatusOK
	pagination := helpers.GeneratePaginationRequest(ctx)

	response, err := c.postService.CommentMost(ctx, pagination)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed pagination comment post", err.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if !response.Success {
		code = http.StatusBadRequest
	}

	res := helpers.BuildResponse(true, "Ok", response)
	ctx.AbortWithStatusJSON(code, res)
}

func (c *postController) LikeMost(ctx *gin.Context) {
	code := http.StatusOK
	pagination := helpers.GeneratePaginationRequest(ctx)

	response, err := c.postService.LikeMost(ctx, pagination)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed pagination like post", err.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if !response.Success {
		code = http.StatusBadRequest
	}

	res := helpers.BuildResponse(true, "Ok", response)
	ctx.AbortWithStatusJSON(code, res)
}

func (c *postController) FindByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	read, err := c.postService.FindByUserID(ctx, userID)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed to read post", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	} else {
		response := helpers.BuildResponse(true, "Readed!", read)
		ctx.JSON(http.StatusOK, response)
	}
}

func (c *postController) getUserIDByToken(token string) string {
	Token, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := Token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
