package controllers

import (
	"github.com/gin-gonic/gin"
	"instatasks/models"
	"log"
	"net/http"
	"strconv"
)

type Task = models.Task

func CreateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			task      Task
			user      User
			userAgent UserAgent
			price     uint
		)

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			panic(err)
			return
		}

		if err := c.ShouldBindHeader(&userAgent); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		if err := userAgent.FindPrice(); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Uncnown User-Agent"})
			log.Println("Error: Uncnown User-Agent")
			return
		}

		user.Instagramid = task.Instagramid
		if err := user.First(); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
			log.Println("Error: User Not Found")
			return
		}

		switch task.Type {
		case "like":
			price = userAgent.Pricelike
		case "follow":
			price = userAgent.Pricefollow
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong Task Type"})
			log.Println("Error: Wrong Task Type")
			return
		}

		total_price := task.Count * price

		if total_price > user.Coins {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Not Enough Coins"}) // 406
			log.Println("Error: Not Enough Coins")
			return
		}

		balance := user.Coins - total_price
		if err := user.UpdateColumn("coins", balance); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err})
			log.Println("Error: ", err)
			return
		}

		task.Create()

		c.JSON(200, gin.H{"coins": balance})
	}
}

func GetHistory() gin.HandlerFunc {
	return func(c *gin.Context) {

		json := struct{ Instagramid uint64 }{}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			panic(err)
			return
		}

		tasks := []models.Task{}

		if err := models.DB.Unscoped().Where("instagramid = ?", json.Instagramid).Order("id desc").Limit(10).Find(&tasks).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err})
			log.Println("Error: ", err)
			return
		}

		var temp []map[string]interface{}

		for _, task := range tasks {
			temp = append(temp, gin.H{
				"taskid":            strconv.FormatUint(uint64(task.ID), 10),
				"created_at":        task.CreatedAt,
				"deleted_at":        task.DeletedAt,
				"type":              task.Type,
				"count":             task.Count,
				"left_counter":      task.LeftCounter,
				"photourl":          task.Photourl,
				"instagramusername": task.Instagramusername,
				"mediaid":           task.Mediaid,
			})
		}

		c.JSON(200, temp)
	}
}
