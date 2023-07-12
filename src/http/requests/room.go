package requests

type MakeRoomRequest struct {
	UserUUID string `json:"uuid"`
}

type JoinRoomRequest struct {
	UserUUID string `json:"uuid"`
	RoomUUID string `json:"room_uuid"`
}

type StateRequest struct {
	UserUUID string `json:"uuid"`
}

type LeaveRoomRequest struct {
	UserUUID string `json:"uuid"`
}
