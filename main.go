package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

var (
	buildType string
	buildVersion string
)

func main() {
	if buildType == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	port, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		port = "8080"
	}

	router := newRouter()
	router.Run(":"+port)
}
