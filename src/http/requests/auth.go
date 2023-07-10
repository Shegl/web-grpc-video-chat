package requests

type AuthRequest struct {
	UserName string `form:"username"`
}

type LogoutRequest struct {
	UUID string `form:"uuid"`
}
