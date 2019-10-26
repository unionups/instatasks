package controllers

import (
	"github.com/gin-gonic/gin"
	// "github.com/imdario/mergo"
	"instatasks/config"
	"instatasks/database"
	. "instatasks/helpers"
	"instatasks/models"
	"net/http"

	"log"
)

type UserAgent = models.UserAgent

type UserAgentData struct {
	Data UserAgent
}

func GetUseragent() gin.HandlerFunc {
	return func(c *gin.Context) {

		var userAgent UserAgent

		db := database.GetDB()

		if err := c.ShouldBindHeader(&userAgent); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		if err := db.First(&userAgent).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Uncnown User-Agent"})
			log.Println("Error: Uncnown User-Agent")
			return
		}

		c.JSON(200, gin.H{
			"activitylimit": userAgent.Activitylimit,
			"like":          userAgent.Like,
			"follow":        userAgent.Follow,
			"pricefollow":   userAgent.Pricefollow,
			"pricelike":     userAgent.Pricelike,
		})
	}
}

func CreateUseragent() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.MustGet(gin.AuthUserKey)

		var json UserAgentData
		var rsaKey models.RsaKey

		db := database.GetDB()

		if err = c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userAgent := json.Data

		if err := db.Create(&userAgent).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err})
			log.Println("Error: Can't create User-Agent")
			return
		}
		db.Save(&userAgent)

		if err := db.Model(&userAgent).Related(&rsaKey, "RsaKey").Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			log.Println("Error: Can't get RSA key")
			return
		}

		rsa_public_key := string(AesDecrypt(rsaKey.RsaPublicKeyAesEncripted, config.GetConfig().Server.AesPassphrase))

		c.JSON(200, gin.H{
			"name":           userAgent.Name,
			"activitylimit":  userAgent.Activitylimit,
			"like":           userAgent.Like,
			"follow":         userAgent.Follow,
			"pricefollow":    userAgent.Pricefollow,
			"pricelike":      userAgent.Pricelike,
			"rsa_public_key": rsa_public_key,
		})
	}
}

func GetRsaPublicKey() gin.HandlerFunc {
	return func(c *gin.Context) {

		var rsaKey models.RsaKey

		db := database.GetDB()

		if err = c.ShouldBindJSON(&rsaKey); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.First(&rsaKey).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err})
			log.Println("Error: RSA key not found")
			return
		}

		rsa_public_key := string(AesDecrypt(rsaKey.RsaPublicKeyAesEncripted, config.GetConfig().Server.AesPassphrase))

		c.JSON(200, gin.H{"rsa_public_key": rsa_public_key})
	}
}
