package controllers

import "github.com/gin-gonic/gin"

type RoomController struct {
}

func (c *RoomController) Make(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func (c *RoomController) Join(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func (c *RoomController) StreamPush(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func (c *RoomController) StreamReceive(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func NewRoomController() *RoomController {
	return &RoomController{}
}
