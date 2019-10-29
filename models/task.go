package models

import (
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	. "instatasks/helpers"
	"instatasks/redis_storage"
	"strconv"
	"time"
)

type Task struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// Taskid            string     `json:"taskid" binding:"-"`
	Type              string `json:"type" binding:"required"`
	Count             uint   `json:"count" binding:"required" gorm:"not null"`
	LeftCounter       uint
	Photourl          string `json:"photourl"`
	Instagramusername string `json:"instagramusername"`
	Mediaid           string `json:"mediaid"`

	Instagramid uint64 `json:"instagramid" binding:"required" sql:"index" gorm:"not null"`

	DoneUsers []*User `gorm:"many2many:user_task"`
}

type CachedTask struct {
	// Taskid            string
	ID                uint
	Type              string
	Count             uint
	LeftCounter       uint
	Photourl          string
	Instagramusername string
	Mediaid           string
}

var (
	TaskRedisCacheCodec *redis_storage.CacheCodec
)

// func (t *Task) AfterCreate(tx *gorm.DB) (err error) {
// 	tx.Model(t).UpdateColumns(Task{
// 		Taskid:      strconv.FormatUint(uint64(t.ID), 10),
// 		LeftCounter: t.Count,
// 	})
// 	return
// }
func (t *Task) BeforeCreate() (err error) {
	t.LeftCounter = t.Count
	return
}

func (t *Task) BeforeUpdate() (err error) {
	t.LeftCounter--
	return
}

func (t *Task) AfterUpdate(tx *gorm.DB) (err error) {
	if t.LeftCounter == 0 {
		tx.Delete(t)
		t.ClearCachedTask()
	}
	return
}

func InitTaskCache() {
	TaskRedisCacheCodec = redis_storage.RegisterCacheCodec("Task")
}

func (t *Task) Create() (err error) {
	if err = DB.Create(t).Error; err != nil {
		return
	}
	t.CacheTask()
	return
}

func (t *Task) DecrementLeftCounter() (err error) {
	err = DB.Model(t).UpdateColumn("left_counter", t.LeftCounter-1).Error
	return
}

func (t *Task) CacheTask() (err error) {
	cachedTask := CachedTask{}
	copier.Copy(&cachedTask, t)

	TaskRedisCacheCodec.Set(&redis_storage.CacheItem{
		Key:        strconv.FormatUint(t.Instagramid, 10),
		Object:     &cachedTask,
		Expiration: DurationInHours(ServerConfig.Cache.NewTaskExpiration),
	})

	return
}

func (t *Task) ClearCachedTask() (err error) {
	TaskRedisCacheCodec.Delete(strconv.FormatUint(t.Instagramid, 10))
	return
}
