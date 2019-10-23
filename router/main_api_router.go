package router

import (
	"github.com/gin-gonic/gin"
	"instatasks/controllers"
	"instatasks/middlwares"
	"net/http"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middlwares.CORS())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/accaunt", controllers.GetOrCreateUser())

	return router
}
