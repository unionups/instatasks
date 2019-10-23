package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/joho/godotenv/autoload"
	"instatasks/config"
	"instatasks/database"
	"instatasks/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var db *gorm.DB
var err error

func main() {

	config := config.InitConfig()
	db = database.InitDB()
	db.DB().Ping()
	defer db.Close()

	if config.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.SetupRouter()

	log.Printf("Listen on port: %s\n", config.Server.Port)

	srv := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: r,
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
