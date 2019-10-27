package models

import (
	. "instatasks/helpers"
	// "github.com/jinzhu/gorm"
	"github.com/jinzhu/copier"
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
	UserRedisCacheCodec *redis_storage.CacheCodec
)

func InitUserCache() {
	UserRedisCacheCodec = redis_storage.RegisterCacheCodec("User")
}

func (user *User) FirstNotBannedOrCreate() (err error) {

	user.FirstOrCreate()

	if user.Banned {
		err = ErrStatusForbidden
		return
	}

	if !DB.First(&BannedDevice{Deviceid: user.Deviceid}).RecordNotFound() {
		err = ErrStatusForbidden
		return
	}
	return
}

func (user *User) FirstOrCreate() (err error) {
	var cachedUser CachedUser

	if err = UserRedisCacheCodec.Once(&redis_storage.CacheItem{
		Key:        user.getIdString(),
		Object:     &cachedUser,
		Expiration: DurationInHours(ServerConfig.Cache.NewUserExpiration),
		Func: func() (interface{}, error) {
			if err = DB.FirstOrCreate(user, &user).Error; err != nil {
				return nil, err
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

	if err = DB.Save(user).Error; err != nil {
		return
	}

	cachedUser := CachedUser{}
	copier.Copy(&cachedUser, user)

	UserRedisCacheCodec.Set(&redis_storage.CacheItem{
		Key:        user.getIdString(),
		Object:     &cachedUser,
		Expiration: DurationInHours(ServerConfig.Cache.UserExpiration),
	})

	return
}

func (user *User) getIdString() string {
	return strconv.FormatUint(user.Instagramid, 10)
}
