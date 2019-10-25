package controllers

import (
	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
	"instatasks/database"
	"instatasks/models"
	"log"
	"net/http"
)

type UAHeader struct {
	Name string `header:"User-Agent" binding:"required"`
}
type UserAgent = models.UserAgent

func ShowUseragent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userAgent *UserAgent
			h         *UAHeader
		)
		db := database.GetDB()

		if err := c.ShouldBindHeader(h); err != nil {
			// if err := c.ShouldBindWith(h, binding.Header); err != nil {
			// type *gin.Context has no field or method ShouldBindHeader
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		userAgent.Name = h.Name

		if err := db.First(userAgent).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Uncnown User-Agent"})
			log.Println("Error: Uncnown User-Agent")
			return
		}

		c.JSON(200, userAgent)
	}
}

func CreateUseragent() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.MustGet(gin.AuthUserKey)
		// user := c.MustGet(gin.AuthUserKey).(string)

		// if secret, ok := secrets[user]; ok {
		// 	c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		// } else {
		// 	c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		// }
		c.JSON(200, nil)
	}
}
