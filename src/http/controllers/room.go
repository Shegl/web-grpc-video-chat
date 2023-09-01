package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"web-grpc-video-chat/src/http/requests"
	services2 "web-grpc-video-chat/src/internal/core/services"
)

type RoomController struct {
	roomService *services2.RoomService
	authService *services2.AuthService
}

func (c *RoomController) Make(ctx *gin.Context) {
	var request requests.MakeRoomRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, nil)
		return
	}
	user, err := c.getUser(request.UserUUID)
	if err != nil {
		ctx.JSON(401, err.Error())
		return
	}
	room, err := c.roomService.Create(user)
	if err != nil {
		ctx.JSON(422, err.Error())
		return
	}
	ctx.JSON(200, room)
}

func (c *RoomController) Join(ctx *gin.Context) {
	var request requests.JoinRoomRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, nil)
		return
	}
	user, err := c.getUser(request.UserUUID)
	if err != nil {
		ctx.JSON(401, err.Error())
		return
	}
	roomUuid, err := uuid.Parse(request.RoomUUID)
	if err != nil {
		ctx.JSON(422, "Cannot parse provided UUID")
		return
	}
	room, err := c.roomService.Join(roomUuid, user)
	if err != nil {
		ctx.JSON(422, err.Error())
		return
	}
	ctx.JSON(200, room)
}

func (c *RoomController) State(ctx *gin.Context) {
	var request requests.StateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, nil)
		return
	}
	user, err := c.getUser(request.UserUUID)
	if err != nil {
		ctx.JSON(401, err.Error())
		return
	}
	room := c.roomService.State(user)
	if room != nil {
		ctx.JSON(200, room)
	} else {
		ctx.JSON(422, nil)
	}
}

func (c *RoomController) Leave(ctx *gin.Context) {
	var request requests.LeaveRoomRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, nil)
		return
	}
	user, err := c.getUser(request.UserUUID)
	if err != nil {
		ctx.JSON(401, err.Error())
		return
	}
	c.roomService.Leave(user)
	ctx.JSON(200, "OK")
}

func (c *RoomController) StreamPush(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func (c *RoomController) StreamReceive(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func (c *RoomController) getUser(userAuth string) (*dto.User, error) {
	userUUID, err := uuid.Parse(userAuth)
	if err != nil {
		return nil, errors.New("Cannot parse provided UUID. ")
	}
	user, err := c.authService.GetUser(userUUID)
	return user, err
}

func NewRoomController(roomService *services2.RoomService, authService *services2.AuthService) *RoomController {
	return &RoomController{roomService: roomService, authService: authService}
}
