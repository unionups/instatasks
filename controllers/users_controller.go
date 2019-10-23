package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"instatasks/database"
	. "instatasks/helpers"
	"instatasks/models"
	"log"
	"net/http"
)

var db *gorm.DB
var err error

type User = models.User

type Data struct {
	Data User
}

func GetOrCreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		db = database.GetDB()
		var json Data

		if err = c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := json.Data

		if user.Deviceid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Request must have Devise ID"})
			log.Println("Bad Request Error: Request must have Devise ID ")
			return
		}

		if err = models.FirstNotBannedUserScope(&user, db); err != nil {
			if gorm.IsRecordNotFoundError(err) {
				// record not found
				if err = db.Create(&user).Error; err != nil {
					c.AbortWithStatusJSON(500, nil)
					log.Println("DB error: ", err)
					return
				}
			} else if IsStatusForbiddenError(err) {
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
