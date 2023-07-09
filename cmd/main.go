package main

import (
	"context"
	"go.uber.org/dig"
	"macos-cam-grpc-chat/src"
)

var Container *dig.Container
var Version = "dev"

func main() {
	Container = src.BuildContainer()
	err := Container.Invoke(func(application *src.Application) {
		ctx := context.Background()
		application.Init(Version)
		application.Run(ctx)
	})

	if err != nil {
		panic(err)
	}
}
