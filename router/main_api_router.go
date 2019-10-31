package router

import (
	"github.com/gin-gonic/gin"
	"instatasks/config"
	"instatasks/controllers"
	"instatasks/middlwares"

	"net/http"
)

var Router *gin.Engine

func SetupRouter() *gin.Engine {
	serverConfig := config.GetConfig().Server

	router := gin.Default()
	router.Use(middlwares.CORS())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/accaunt", controllers.GetOrCreateUser())
	router.POST("/setting", controllers.GetUseragent())
	router.POST("/newwork", controllers.CreateTask())
	router.POST("/history", controllers.GetByUserCreatedTaskHistory())
	router.POST("/gettasks", controllers.GetTasks())
	router.POST("/done", controllers.DoneTask())

	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		serverConfig.Superadmin.Username: serverConfig.Superadmin.Password,
	}))

	authorized.POST("/useragent", controllers.CreateUseragent())
	authorized.GET("/useragent/pkey", controllers.GetRsaPublicKey())

	Router = router

	return router
}

func GetRouter() *gin.Engine {
	return Router
}
