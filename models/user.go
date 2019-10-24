package models

import (
	"github.com/go-redis/cache/v7"
	// "github.com/jinzhu/gorm"
	"github.com/jinzhu/copier"
	"instatasks/config"
	"instatasks/database"
	. "instatasks/helpers"
	"instatasks/redis_storage"
	"strconv"
	"time"
)

type User struct {
	Instagramid uint64 `json:"instagramid" binding:"required" gorm:"primary_key" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Banned      bool       `gorm:"default:false"`
	Coins       int        `json:"coins" gorm:"default:0"`
	Deviceid    string     `json:"deviceid" gorm:"-"`
	Rateus      bool       `json:"rateus" gorm:"default:true"`
}

type CachedUser struct {
	Banned bool
	Coins  int
	Rateus bool
}

func FirstNotBannedOrCreateUserScope(user *User) error {

	FirstOrCreateUser(user)

	if user.Banned {
		return ErrStatusForbidden
	}

	db := database.GetDB()
	if !db.First(&BannedDevice{Deviceid: user.Deviceid}).RecordNotFound() {
		return ErrStatusForbidden
	}
	return nil
}

func FirstOrCreateUser(user *User) error {
	var cachedUser CachedUser

	redis_cache := redis_storage.GetRedisCache()

	if err := redis_cache.Once(&cache.Item{
		Key:        getIdString(user),
		Object:     &cachedUser,
		Expiration: DurationInHours(config.GetConfig().Server.Cache.NewUserExpiration),
		Func: func() (interface{}, error) {
			db := database.GetDB()
			if db.First(user).RecordNotFound() {
				if err := db.Create(user).Error; err != nil {
					return nil, err
				}
			}
			copier.Copy(&cachedUser, user)
			return cachedUser, nil
		},
	}); err != nil {
		return err
	}
	copier.Copy(user, &cachedUser)
	return nil
}

func SaveUser(user *User) error {

	db := database.GetDB()
	redis_cache := redis_storage.GetRedisCache()

	if err := db.Save(&user).Error; err != nil {
		return err
	}

	cachedUser := CachedUser{}
	copier.Copy(&cachedUser, user)

	redis_cache.Set(&cache.Item{
		Key:        getIdString(user),
		Object:     &cachedUser,
		Expiration: DurationInHours(config.GetConfig().Server.Cache.UserExpiration),
	})

	return nil
}

func getIdString(user *User) string {
	return strconv.FormatUint(user.Instagramid, 10)
}
