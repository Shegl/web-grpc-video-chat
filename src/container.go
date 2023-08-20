package src

import (
	"go.uber.org/dig"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/http"
	"web-grpc-video-chat/src/http/controllers"
	"web-grpc-video-chat/src/http/middleware"
	"web-grpc-video-chat/src/inroom"
	"web-grpc-video-chat/src/services"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	processError(container.Provide(NewApplication))
	processError(container.Provide(http.NewWebServer))

	processError(container.Provide(controllers.NewAuthController))
	processError(container.Provide(controllers.NewRoomController))

	processError(container.Provide(dto.NewRepository))

	processError(container.Provide(middleware.NewCorsMiddleware))

	processError(container.Provide(inroom.NewRoomProvider))
	processError(container.Provide(services.NewAuthService))
	processError(container.Provide(services.NewRoomService))
	processError(container.Provide(inroom.NewChatServer))
	processError(container.Provide(inroom.NewStreamServer))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
