package main

import (
	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine{
	router := gin.Default()

	for _, route := range routes {
		router.Handle(route.method, route.path, route.handler)
	}

	return router
}

type route struct {
	method string
	path string
	handler gin.HandlerFunc
}

var routes = []route{
	{
		"GET",
		"/api/download/Mb/:size",
		MbDownloadHandler,
	},
	{
		"POST",
		"/api/upload/Mb/:size",
		MbUploadHandler,
	},
	{
		"GET",
		"/api/ping",
		PingHandler,
	},
}
