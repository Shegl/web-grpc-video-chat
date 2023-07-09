package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type WebServer struct {
	router       *gin.Engine
	srv          *http.Server
	wg           *sync.WaitGroup
	shutdownChan chan struct{}
}

func (w *WebServer) Init(addr string, wg *sync.WaitGroup, shutdownChan chan struct{}) error {
	w.shutdownChan = shutdownChan
	w.wg = wg

	return nil
}

func (w *WebServer) Run(ctx context.Context) error {
	return nil
}

func NewWebServer() *WebServer {
	webServer := &WebServer{router: nil}
	return webServer
}
