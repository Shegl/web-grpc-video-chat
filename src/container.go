package src

import (
	"go.uber.org/dig"
	"web-grpc-video-chat/src/chat"
	"web-grpc-video-chat/src/http"
	"web-grpc-video-chat/src/http/controllers"
	"web-grpc-video-chat/src/http/middleware"
	"web-grpc-video-chat/src/services"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	processError(container.Provide(NewApplication))
	processError(container.Provide(http.NewWebServer))
	processError(container.Provide(chat.NewChatServiceServer))

	processError(container.Provide(controllers.NewAuthController))
	processError(container.Provide(controllers.NewRoomController))

	processError(container.Provide(middleware.NewCorsMiddleware))

	processError(container.Provide(services.NewAuthService))
	processError(container.Provide(services.NewRoomService))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
