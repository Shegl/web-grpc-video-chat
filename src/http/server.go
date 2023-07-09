package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"macos-cam-grpc-chat/src/http/controllers"
	"net/http"
	"sync"
	"time"
)

type WebServer struct {
	router         *gin.Engine
	srv            *http.Server
	wg             *sync.WaitGroup
	authController *controllers.AuthController
	roomController *controllers.RoomController
	shutdownChan   chan struct{}
}

func (w *WebServer) Init(addr string, wg *sync.WaitGroup, shutdownChan chan struct{}) error {
	w.shutdownChan = shutdownChan
	w.wg = wg
	w.router = gin.New()
	w.initRoutes()

	w.srv = &http.Server{
		Addr:    addr,
		Handler: w.router,
	}

	return nil
}

func (w *WebServer) Run(ctx context.Context) error {
	httpShutdownCh := make(chan struct{})

	go func() {
		<-ctx.Done()

		log.Println("WebServer:: shutdown started")

		graceCtx, graceCancel := context.WithTimeout(ctx, 1*time.Second)
		defer graceCancel()

		if err := w.srv.Shutdown(graceCtx); err != nil {
			log.Println(err)
		}

		httpShutdownCh <- struct{}{}
	}()

	w.wg.Add(1)

	go func() {
		defer w.wg.Done()

		err := w.srv.ListenAndServe()

		<-httpShutdownCh

		if err != http.ErrServerClosed {
			panic(err)
		}

		log.Println("WebServer:: shutdown complete")
		close(w.shutdownChan)
	}()

	return nil
}

func (w *WebServer) initRoutes() {
	w.router.POST("/auth", w.authController.Auth)

	w.router.POST("/make-room", w.roomController.Make)
	w.router.POST("/join-room", w.roomController.Join)

	w.router.POST("/room/:id", w.roomController.StreamPush)
	w.router.GET("/room/:id", w.roomController.StreamReceive)
}

func NewWebServer(
	authController *controllers.AuthController,
	roomController *controllers.RoomController,
) *WebServer {
	webServer := &WebServer{
		authController: authController,
		roomController: roomController,

		router: nil,
	}
	return webServer
}
