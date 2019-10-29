package controllers

import (
	"github.com/gin-gonic/gin"
	. "instatasks/helpers"
	"instatasks/models"
	"net/http"

	"log"
)

type User = models.User

func GetOrCreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			panic(err)
			return
		}

		if user.Deviceid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Request must have Devise ID"})
			log.Println("Bad Request Error: Request must have Devise ID ")
			return
		}

		if err := user.FirstNotBannedOrCreate(); err != nil {
			if IsStatusForbiddenError(err) {
				c.AbortWithStatusJSON(http.StatusForbidden, nil)
				log.Println("Error: Banned")
				return
			} else {
				c.AbortWithStatusJSON(500, nil)
				log.Println("DB error: ", err)
				return
			}
		}
		c.JSON(200, gin.H{"instagramid": user.Instagramid, "coins": user.Coins, "rateus": user.Rateus})
	}
}
