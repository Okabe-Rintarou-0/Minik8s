package apiserver

import (
	"github.com/gin-gonic/gin"
	"log"
)

type ApiServer interface {
	Run()
}

func New() ApiServer {
	return &apiServer{
		httpServer: gin.Default(),
	}
}

type apiServer struct {
	httpServer *gin.Engine
}

func (api *apiServer) bindHandlers() {
	for url, handler := range postTable {
		api.httpServer.POST(url, handler)
	}

	for url, handler := range getTable {
		api.httpServer.GET(url, handler)
	}

	for url, handler := range deleteTable {
		api.httpServer.DELETE(url, handler)
	}

	for url, handler := range putTable {
		api.httpServer.PUT(url, handler)
	}
}

func (api *apiServer) Run() {
	api.bindHandlers()
	log.Fatal(api.httpServer.Run())
}
