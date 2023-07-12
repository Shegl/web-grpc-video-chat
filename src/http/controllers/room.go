package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"macos-cam-grpc-chat/src/dto"
	"macos-cam-grpc-chat/src/http/requests"
	"macos-cam-grpc-chat/src/services"
)

type RoomController struct {
	roomService *services.RoomService
	authService *services.AuthService
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
	roomUUID, err := uuid.Parse(request.RoomUUID)
	if err != nil {
		ctx.JSON(422, "Cannot parse provided UUID")
		return
	}
	room, err := c.roomService.Join(roomUUID, user)
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

func NewRoomController(roomService *services.RoomService, authService *services.AuthService) *RoomController {
	return &RoomController{roomService: roomService, authService: authService}
}
