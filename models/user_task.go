package models

import (
	"instatasks/redis_storage"
)

type CachedUserTask struct {
	UserInstagramid uint64
	TaskId          uint
}

var (
	UserTaskRedisCacheCodec *redis_storage.CacheCodec
)

func InitUserTaskCache() {
	TaskRedisCacheCodec = redis_storage.RegisterCacheCodec("UserTask")
}
