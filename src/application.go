package src

import (
	"context"
	"fmt"
	"log"
	"macos-cam-grpc-chat/src/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Application struct {
	webServer    *http.WebServer
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

	a.webServer.Init(
		":3000",
		&a.wg,
		a.shutdownChan,
	)

	fmt.Println("application:: Init() :: init complete")
	return nil
}

func (a *Application) Run(ctx context.Context) {
	fmt.Println("application:: Run() :: starting")
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	a.processSignals(cancelFunc)

	err := a.webServer.Run(cancelCtx)
	if err != nil {
		panic(err)
	}

	fmt.Println("application:: Run() :: running")
	a.wg.Wait()

	fmt.Println("application:: Run() :: graceful shutdown")
}

func (a *Application) processSignals(cancelFunc context.CancelFunc) {
	go func() {
		select {
		case <-a.sigs:
			log.Println("application: received shutdown signal from OS")
			cancelFunc()
			break
		}
	}()
}

func NewApplication(
	webServer *http.WebServer,
) *Application {
	return &Application{
		webServer: webServer,
	}
}
