package src

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"web-grpc-video-chat/src/http"
	"web-grpc-video-chat/src/inroom"
)

type Application struct {
	webServer    *http.WebServer
	chatServer   *inroom.ChatServer
	streamServer *inroom.StreamServer
	wg           sync.WaitGroup
	sigs         chan os.Signal
	shutdownChan chan struct{}

	Version string
}

func (a *Application) Init(version string) error {
	a.Version = version

	// for graceful shutdown
	a.sigs = make(chan os.Signal, 1)
	a.shutdownChan = make(chan struct{}, 1)

	signal.Notify(a.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	a.webServer.Init(":3000", &a.wg, a.shutdownChan)
	a.chatServer.Init(":3001", &a.wg)
	a.streamServer.Init(":3002", &a.wg)

	log.Println("application:: Init() :: init complete")
	return nil
}

func (a *Application) Run(ctx context.Context) {
	log.Println("application:: Run() :: starting")
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	a.processSignals(cancelFunc)

	err := a.webServer.Run(cancelCtx)
	if err != nil {
		panic(err)
	}

	err = a.chatServer.Run(cancelCtx)
	if err != nil {
		panic(err)
	}

	err = a.streamServer.Run(cancelCtx)
	if err != nil {
		panic(err)
	}

	log.Println("application:: Run() :: running")
	a.wg.Wait()

	log.Println("application:: Run() :: graceful shutdown")
}

func (a *Application) processSignals(cancelFunc context.CancelFunc) {
	go func() {
		select {
		case <-a.sigs:
			log.Println("application:: received shutdown signal from OS")
			cancelFunc()
			break
		}
	}()
}

func NewApplication(
	webServer *http.WebServer,
	chatServer *inroom.ChatServer,
	streamServer *inroom.StreamServer,
) *Application {
	return &Application{
		webServer:    webServer,
		chatServer:   chatServer,
		streamServer: streamServer,
	}
}
