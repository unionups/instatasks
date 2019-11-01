package router

import (
	"github.com/aviddiviner/gin-limit"
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

	if serverConfig.ConnectionLimit > 0 {
		router.Use(limit.MaxAllowed(serverConfig.ConnectionLimit))
	}

	router.Use(middlwares.CORS())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	mainApi := router.Group("/")

	if serverConfig.BodyCrypt {
		mainApi.Use(middlwares.BodyCrypt())
	}

	mainApi.POST("/accaunt", controllers.GetOrCreateUser())
	mainApi.POST("/setting", controllers.GetUseragent())
	mainApi.POST("/newwork", controllers.CreateTask())
	mainApi.POST("/history", controllers.GetByUserCreatedTaskHistory())
	mainApi.POST("/gettasks", controllers.GetTasks())
	mainApi.POST("/done", controllers.DoneTask())
	mainApi.POST("/rateus", controllers.DoneRateus())

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
