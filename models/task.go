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
	ID                uint       `json:"-" gorm:"primary_key"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at" sql:"index"`
	Taskid            string     `json:"taskid" gorm:"-"`
	Type              string     `json:"type" binding:"required"`
	Count             uint       `json:"count" binding:"required" gorm:"not null"`
	LeftCounter       uint       `json:"left_counter"`
	Photourl          string     `json:"photourl"`
	Instagramusername string     `json:"instagramusername"`
	Mediaid           string     `json:"mediaid" binding:"required" sql:"index" gorm:"not null"`
	CancelLeftCounter uint8      `json:"-"`

	Instagramid uint `json:"instagramid" binding:"required" sql:"index" gorm:"not null"`
}

type CachedTask struct {
	Type              string
	LeftCounter       uint
	CancelLeftCounter uint8
	Mediaid           string
}

var (
	TaskRedisCacheCodec *redis_storage.CacheCodec
)

func (t *Task) BeforeCreate() (err error) {
	t.LeftCounter = t.Count
	t.CancelLeftCounter = 10
	return
}

func (t *Task) AfterUpdate(tx *gorm.DB) (err error) {
	if (t.LeftCounter == 0) || (t.CancelLeftCounter == 0) {
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

func (t *Task) Done() (err error) {
	err = t.DecrementLeftCounter()
	return
}

func (t *Task) Cancel() (err error) {
	err = t.DecrementCancelLeftCounter()
	return
}

func (t *Task) DecrementLeftCounter() (err error) {

	cachedTask := CachedTask{}

	if err = TaskRedisCacheCodec.Get(t.getIdString(), &cachedTask); err == nil {

		cachedTask.LeftCounter--
		copier.Copy(t, &cachedTask)

		go DB.Model(t).Update("left_counter", t.LeftCounter)

		err = TaskRedisCacheCodec.Set(&redis_storage.CacheItem{
			Key:        t.getIdString(),
			Object:     &cachedTask,
			Expiration: DurationInHours(ServerConfig.Cache.NewTaskExpiration),
		})
		return
	}

	conn, tx := NewTransaction(DB.Begin())

	if err = conn.Select("type, left_counter, mediaid, cancel_left_counter").First(t).Error; err != nil {
		tx.Fail()
	}
	if err = conn.Model(t).Update("left_counter", t.LeftCounter-1).Error; err != nil {
		tx.Fail()
	}
	tx.Close()

	t.CacheTask()
	return
}

func (t *Task) DecrementCancelLeftCounter() (err error) {

	cachedTask := CachedTask{}

	if err = TaskRedisCacheCodec.Get(t.getIdString(), &cachedTask); err == nil {

		cachedTask.CancelLeftCounter--
		copier.Copy(t, &cachedTask)

		go DB.Model(t).Update("cancel_left_counter", t.CancelLeftCounter)

		err = TaskRedisCacheCodec.Set(&redis_storage.CacheItem{
			Key:        t.getIdString(),
			Object:     &cachedTask,
			Expiration: DurationInHours(ServerConfig.Cache.TaskExpiration),
		})
		return
	}

	conn, tx := NewTransaction(DB.Begin())

	if err = conn.Select("type, left_counter, mediaid, cancel_left_counter").First(t).Error; err != nil {
		tx.Fail()
	}
	if err = conn.Model(t).Update("cancel_left_counter", t.CancelLeftCounter-1).Error; err != nil {
		tx.Fail()
	}
	tx.Close()
	t.CacheTask()
	return
}

func (t *Task) CacheTask() (err error) {
	cachedTask := CachedTask{}
	copier.Copy(&cachedTask, t)

	TaskRedisCacheCodec.Set(&redis_storage.CacheItem{
		Key:        t.getIdString(),
		Object:     &cachedTask,
		Expiration: DurationInHours(ServerConfig.Cache.NewTaskExpiration),
	})

	return
}

func (t *Task) ClearCachedTask() (err error) {
	TaskRedisCacheCodec.Delete(t.getIdString())
	return
}

func (t *Task) getIdString() string {
	return strconv.FormatUint(uint64(t.ID), 10)
}
