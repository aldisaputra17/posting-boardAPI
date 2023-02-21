package controllers

import (
	"fmt"
	"net/http"

	"github.com/aldisaputra17/posting-board/dto"
	"github.com/aldisaputra17/posting-board/entities"
	"github.com/aldisaputra17/posting-board/helpers"
	"github.com/aldisaputra17/posting-board/services"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func NewAuthController(authService services.AuthService, jwtService services.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	var reqRegister *dto.RequestRegister
	errObj := ctx.ShouldBind(&reqRegister)
	if errObj != nil {
		response := helpers.BuildErrorResponse("Failed to process request", errObj.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(reqRegister.Email) {
		response := helpers.BuildErrorResponse("Failed to process request", "Duplicate username", helpers.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
		return
	}
	createdUser, err := c.authService.CreateUser(ctx, reqRegister)
	if err != nil {
		res := helpers.BuildErrorResponse("Failed to created", err.Error(), helpers.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		fmt.Println("erorr", err)
		return
	} else {
		token := c.jwtService.GenerateToken(createdUser.ID)
		createdUser.Token = token
		response := helpers.BuildResponse(true, "Created!", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var reqLogin dto.RequestLogin
	err := ctx.ShouldBind(&reqLogin)
	if err != nil {
		response := helpers.BuildErrorResponse("Failed to process request", err.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(reqLogin.Email, reqLogin.Password)
	if v, ok := authResult.(entities.User); ok {
		generatedToken := c.jwtService.GenerateToken(v.ID)
		v.Token = generatedToken
		response := helpers.BuildResponse(true, "Ok!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helpers.BuildErrorResponse("Please check again your credential", "Invalid Credential", helpers.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}
