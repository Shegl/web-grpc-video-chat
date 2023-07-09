package src

import (
	"go.uber.org/dig"
	"macos-cam-grpc-chat/src/http"
	"macos-cam-grpc-chat/src/http/controllers"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	processError(container.Provide(NewApplication))
	processError(container.Provide(http.NewWebServer))

	processError(container.Provide(controllers.NewAuthController))
	processError(container.Provide(controllers.NewRoomController))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
