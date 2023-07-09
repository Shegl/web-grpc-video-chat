package src

import (
	"go.uber.org/dig"
	"macos-cam-grpc-chat/src/http"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	processError(container.Provide(NewApplication))
	processError(container.Provide(http.NewWebServer))
	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
