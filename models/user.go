package models

import (
	"github.com/jinzhu/copier"
	// "github.com/jinzhu/gorm"
	. "instatasks/helpers"
	"instatasks/redis_storage"
	"strconv"
	"time"
)

type User struct {
	Instagramid uint `json:"instagramid" binding:"required" gorm:"auto_increment:false;primary_key:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Banned      bool       `gorm:"default:false"`
	Coins       uint       `json:"coins" gorm:"default:0"`
	Deviceid    string     `json:"deviceid" gorm:"-"`
	Rateus      bool       `binding:"-" gorm:"default:true"`

	Tasks []Task `gorm:"foreignkey:Instagramid;association_foreignkey:Instagramid;"`
}

type CachedUser struct {
	Banned bool
	Coins  uint
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
		user.UpdateColumn("banned", true)
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
	copier.Copy(user, &cachedUser)
	return
}

func (user *User) First() (err error) {
	var cachedUser CachedUser

	if err = UserRedisCacheCodec.Once(&redis_storage.CacheItem{
		Key:        user.getIdString(),
		Object:     &cachedUser,
		Expiration: DurationInHours(ServerConfig.Cache.UserExpiration),
		Func: func() (interface{}, error) {
			if err = DB.First(user).Error; err != nil {
				return nil, err
			}
			copier.Copy(&cachedUser, user)
			return cachedUser, nil
		},
	}); err != nil {
		return
	}
	copier.Copy(user, &cachedUser)
	return
}

func (user *User) Save() (err error) {
	if err = DB.Save(user).Error; err != nil {
		return
	}
	err = user.SetCache(ServerConfig.Cache.UserExpiration)
	return
}

func (user *User) UpdateColumns(p map[string]interface{}) (err error) {
	if err = DB.Model(user).UpdateColumns(p).Error; err != nil {
		return
	}
	err = user.SetCache(ServerConfig.Cache.UserExpiration)
	return
}

func (user *User) Updates(p map[string]interface{}) (err error) {
	if err = DB.Model(user).Updates(p).Error; err != nil {
		return
	}
	err = user.SetCache(ServerConfig.Cache.UserExpiration)
	return
}

func (user *User) UpdateColumn(values ...interface{}) (err error) {
	if err = DB.Model(user).UpdateColumn(values).Error; err != nil {
		return
	}
	err = user.SetCache(ServerConfig.Cache.UserExpiration)
	return
}

func (user *User) Update(attrs ...interface{}) (err error) {
	if err = DB.Model(user).Update(attrs).Error; err != nil {
		return
	}
	err = user.SetCache(ServerConfig.Cache.UserExpiration)
	return
}

func (user *User) SetCache(expiration int) (err error) {
	cachedUser := CachedUser{}
	copier.Copy(&cachedUser, user)

	err = UserRedisCacheCodec.Set(&redis_storage.CacheItem{
		Key:        user.getIdString(),
		Object:     &cachedUser,
		Expiration: DurationInHours(expiration),
	})

	return
}

func (user *User) getIdString() string {
	return strconv.FormatUint(uint64(user.Instagramid), 10)
}
