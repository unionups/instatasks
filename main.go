package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/joho/godotenv/autoload"
	"instatasks/config"
	"instatasks/controllers"
	"instatasks/database"
	"instatasks/middlwares"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var db *gorm.DB
var err error

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

func main() {

	config := config.InitConfig()
	db = database.InitDB()
	db.DB().Ping()
	defer db.Close()

	database.Migrate()

	if config.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := SetupRouter()

	log.Printf("Listen on port: %s\n", config.Server.Port)
	srv := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
