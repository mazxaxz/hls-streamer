package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/mazxaxz/hls-streamer/pkg/rest"
)

func setupRouting(handlers ...rest.SetupRouterer) http.Handler {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	v1 := router.Group("v1")
	for _, handler := range handlers {
		handler.SetupRouter(v1)
	}
	return router
}
