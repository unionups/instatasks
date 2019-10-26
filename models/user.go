package models

import (
	. "instatasks/helpers"
	// "github.com/jinzhu/gorm"
	"github.com/jinzhu/copier"
	"instatasks/config"
	"instatasks/database"
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
	Rateus      bool       `binding:"-" gorm:"default:true"`
}

type CachedUser struct {
	Banned bool
	Coins  int
	Rateus bool
}

var (
	RedisCacheCodec *redis_storage.CacheCodec
)

func InitUserCache() {
	RedisCacheCodec = redis_storage.RegisterCacheCodec("User")
}

func (user *User) FirstNotBannedOrCreate() (err error) {

	user.FirstOrCreate()

	if user.Banned {
		err = ErrStatusForbidden
		return
	}

	db := database.GetDB()
	if !db.First(&BannedDevice{Deviceid: user.Deviceid}).RecordNotFound() {
		err = ErrStatusForbidden
		return
	}
	return
}

func (user *User) FirstOrCreate() (err error) {
	var cachedUser CachedUser

	if err = RedisCacheCodec.Once(&redis_storage.CacheItem{
		Key:        user.getIdString(),
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
		return
	}
	copier.Copy(&user, &cachedUser)
	return
}

func (user *User) Save() (err error) {

	db := database.GetDB()

	if err = db.Save(user).Error; err != nil {
		return
	}

	cachedUser := CachedUser{}
	copier.Copy(&cachedUser, user)

	RedisCacheCodec.Set(&redis_storage.CacheItem{
		Key:        user.getIdString(),
		Object:     &cachedUser,
		Expiration: DurationInHours(config.GetConfig().Server.Cache.UserExpiration),
	})

	return
}

func (user *User) getIdString() string {
	return strconv.FormatUint(user.Instagramid, 10)
}
