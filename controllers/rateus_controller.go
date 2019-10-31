package controllers

import (
	"github.com/gin-gonic/gin"
	"instatasks/models"

	"log"
	"net/http"
)

func DoneRateus() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user      User
			userAgent UserAgent
		)

		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
			log.Println("Error: User Not Found")
			return
		}

		if err := c.ShouldBindHeader(&userAgent); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		user.First()

		if !user.Rateus {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Task already done"}) // 406
			log.Println("Error: Task already done")
			return
		}
		if err := models.DB.First(&userAgent).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		user.Coins += userAgent.Pricerateus
		user.Rateus = false

		if err := user.Save(); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		c.JSON(200, gin.H{"coins": user.Coins})

	}
}
