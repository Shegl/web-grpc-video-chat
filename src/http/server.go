package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"macos-cam-grpc-chat/src/http/controllers"
	"macos-cam-grpc-chat/src/http/middleware"
	"net/http"
	"sync"
	"time"
)

type WebServer struct {
	router *gin.Engine
	srv    *http.Server
	wg     *sync.WaitGroup

	authController *controllers.AuthController
	roomController *controllers.RoomController

	corsMiddleware *middleware.CorsMiddleware

	shutdownChan chan struct{}
}

func (w *WebServer) Init(addr string, wg *sync.WaitGroup, shutdownChan chan struct{}) error {
	w.shutdownChan = shutdownChan
	w.wg = wg

	w.router = gin.New()
	w.registerGlobalMiddlewares()
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
	w.router.POST("/check", w.authController.Check)
	w.router.POST("/logout", w.authController.Logout)

	w.router.POST("/room/make", w.roomController.Make)
	w.router.POST("/room/join", w.roomController.Join)
	w.router.POST("/room/state", w.roomController.State)
	w.router.POST("/room/leave", w.roomController.Leave)

	w.router.POST("/room/:id", w.roomController.StreamPush)
	w.router.GET("/room/:id", w.roomController.StreamReceive)
}

func (w *WebServer) registerGlobalMiddlewares() {
	w.router.Use(w.corsMiddleware.Handle)
}

func NewWebServer(
	authController *controllers.AuthController,
	roomController *controllers.RoomController,
	corsMiddleware *middleware.CorsMiddleware,
) *WebServer {
	webServer := &WebServer{
		authController: authController,
		roomController: roomController,
		corsMiddleware: corsMiddleware,

		router: nil,
	}
	return webServer
}
