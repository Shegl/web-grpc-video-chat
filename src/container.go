package src

import (
	"go.uber.org/dig"
	"web-grpc-video-chat/src/http"
	"web-grpc-video-chat/src/http/controllers"
	"web-grpc-video-chat/src/http/middleware"
	services2 "web-grpc-video-chat/src/internal/core/services"
	"web-grpc-video-chat/src/pb"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	processError(container.Provide(NewApplication))
	processError(container.Provide(http.NewWebServer))

	processError(container.Provide(controllers.NewAuthController))
	processError(container.Provide(controllers.NewRoomController))

	processError(container.Provide(services2.NewRoomProvider))
	processError(container.Provide(pb.NewChatServer))
	processError(container.Provide(pb.NewStreamServer))

	processError(container.Provide(services2.NewAuthService))
	processError(container.Provide(services2.NewRoomService))

	processError(container.Provide(repo.NewRepository))

	processError(container.Provide(middleware.NewCorsMiddleware))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
