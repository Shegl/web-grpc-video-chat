package controllers

import (
	"github.com/gin-gonic/gin"
)

type AuthController struct {
}

func (c *AuthController) Auth(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func NewAuthController() *AuthController {
	return &AuthController{}
}
