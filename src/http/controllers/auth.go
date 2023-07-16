package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/http/requests"
	"web-grpc-video-chat/src/services"
)

type AuthController struct {
	authService *services.AuthService
}

func (c *AuthController) Check(ctx *gin.Context) {
	var request requests.CheckRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, nil)
		return
	}
	var user *dto.User
	userUUID, err := uuid.Parse(request.UUID)
	if err == nil {
		user, err = c.authService.GetUser(userUUID)
		if err != nil {
			ctx.JSON(401, nil)
			return
		}
		ctx.JSON(200, user)
	} else {
		ctx.JSON(401, nil)
	}
}

func (c *AuthController) Auth(ctx *gin.Context) {
	var request requests.AuthRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, nil)
		return
	}
	user, err := c.authService.Authenticate(request.UserName)
	if err != nil {
		ctx.JSON(422, err)
	}
	ctx.JSON(200, user)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	var request requests.LogoutRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, nil)
		return
	}

	userUUID, err := uuid.Parse(request.UUID)
	if err == nil {
		c.authService.Logout(userUUID)
	}
	ctx.JSON(200, "OK")
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}
