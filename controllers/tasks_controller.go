package controllers

import (
	"github.com/gin-gonic/gin"
	"instatasks/models"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Task = models.Task

var wg sync.WaitGroup

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
			log.Println("Error: ", err.Error())
			return
		}

		if err := c.ShouldBindHeader(&userAgent); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		user.Instagramid = task.Instagramid
		if err := user.First(); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
			log.Println("Error: User Not Found")
			return
		}

		if err := userAgent.FindPrice(); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Uncnown User-Agent"})
			log.Println("Error: Uncnown User-Agent")
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

		if err := task.Create(); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		c.JSON(200, gin.H{"coins": balance})
	}
}

func GetByUserCreatedTaskHistory() gin.HandlerFunc {
	return func(c *gin.Context) {

		json := struct{ Instagramid uint }{}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			panic(err)
			return
		}

		tasks := []Task{}

		if err := models.DB.Unscoped().Where("instagramid = ?", json.Instagramid).Order("id desc").Limit(10).Find(&tasks).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		for idx, task := range tasks {
			tasks[idx].Taskid = strconv.FormatUint(uint64(task.ID), 10)
		}

		c.JSON(200, tasks)
	}
}

func GetTasks() gin.HandlerFunc {
	return func(c *gin.Context) {

		json := struct {
			Instagramid uint
			Type        string
		}{}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		tasks := []struct {
			Id                string `json:"taskid"`
			Type              string `json:"type"`
			Photourl          string `json:"photourl"`
			Instagramusername string `json:"instagramusername" sql:"index"`
			Mediaid           string `json:"mediaid" sql:"index"`
		}{}

		if json.Type == "all" {
			if err := models.DB.Raw(`SELECT DISTINCT "tasks"."id", "tasks"."type", "tasks"."photourl", "tasks"."instagramusername", "tasks"."mediaid"
				FROM "tasks" 
				LEFT OUTER JOIN "user_mediaids" AS a ON "tasks"."mediaid" = a."mediaid"
				LEFT OUTER JOIN "user_mediaids" AS b ON b.instagramid = $1::integer
				WHERE "tasks"."instagramid" <> $1::integer
				AND "tasks"."deleted_at" IS NULL 
				LIMIT 5;`, json.Instagramid).
				Scan(&tasks).Error; err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				log.Println("Error: ", err.Error())
				return
			}
		} else {
			if err := models.DB.Raw(`SELECT DISTINCT "tasks"."id", "tasks"."type", "tasks"."photourl", "tasks"."instagramusername", "tasks"."mediaid"
				FROM "tasks" 
				LEFT OUTER JOIN "user_mediaids" AS a ON "tasks"."mediaid" = a."mediaid"
				LEFT OUTER JOIN "user_mediaids" AS b ON b.instagramid = $1::integer
				WHERE "tasks"."type" = $2::text
				AND "tasks"."instagramid" <> $1::integer
				AND "tasks"."deleted_at" IS NULL 
				LIMIT 5;`, json.Instagramid, json.Type).
				Scan(&tasks).Error; err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				log.Println("Error: ", err.Error())
				return
			}
		}

		c.JSON(200, tasks)
	}
}
func DoneTask() gin.HandlerFunc {
	return func(c *gin.Context) {

		json := struct {
			Instagramid uint   `binding:"required"`
			Taskid      string `binding:"required"`
			Status      string `binding:"required"`
		}{}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}

		tid, _ := strconv.ParseUint(json.Taskid, 10, 64)

		task := Task{ID: uint(tid)}

		if json.Status == "cancel" {
			/// canceled
			if err := task.Cancel(); err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				log.Println("Error: ", err.Error())
				return
			}
			c.AbortWithStatusJSON(http.StatusResetContent, gin.H{"error": "Task Canceled"}) // 205
			models.DB.Create(&models.UserMediaid{Instagramid: json.Instagramid, Mediaid: task.Mediaid})
			return
		}

		if err := task.Done(); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			log.Println("Error: ", err.Error())
			return
		}
		//// зачислить монеты
		userAgent := UserAgent{}
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

		user := User{Instagramid: json.Instagramid}
		if err := user.First(); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
			log.Println("Error: User Not Found")
			return
		}

		switch task.Type {
		case "like":
			user.Coins += userAgent.Pricelike
		case "follow":
			user.Coins += userAgent.Pricefollow
		}

		c.JSON(200, gin.H{"coins": user.Coins})
		wg.Add(2)
		go func() {
			models.DB.Create(&models.UserMediaid{Instagramid: json.Instagramid, Mediaid: task.Mediaid})
			wg.Done()
		}()
		go func() {
			user.UpdateColumn("coins", user.Coins)
			wg.Done()
		}()
		wg.Wait()
	}
}
