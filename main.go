package main

import (
	"log"
	"main/pods"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Example API
// @version 1.0
// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	// Gin instance
	r := gin.New()

	url := ginSwagger.URL("http://localhost:3000/swagger/docs.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.POST("/pod/:podName/:image", pods.CreatePod)
	r.GET("/expose/:podName/:port", pods.ExposePod)

	// Start server
	if err := r.Run("localhost:3000"); err != nil {
		log.Fatal(err)
	}
}
