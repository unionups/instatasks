package models

import (
	"github.com/jinzhu/copier"
	. "instatasks/helpers"
	"sync"
	"time"
)

type UserAgent struct {
	Name      string `header:"User-Agent" json:"name" binding:"required"  gorm:"primary_key:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Activitylimit uint `json:"activitylimit" gorm:"default:0"`
	Like          bool `json:"like" gorm:"default:true"`
	Follow        bool `json:"follow" gorm:"default:true"`
	Pricefollow   uint `json:"pricefollow" gorm:"default:5"`
	Pricelike     uint `json:"pricelike" gorm:"default:1"`

	RsaKey RsaKey `gorm:"foreignkey:Name;association_foreignkey:Name"`
}

type CachedUserAgentSettings struct {
	Activitylimit uint
	Like          bool
	Follow        bool
	Pricefollow   uint
	Pricelike     uint
}

type CachedPrice struct {
	Pricefollow uint
	Pricelike   uint
}

var (
	cachedUserAgentSettings = make(map[string]*CachedUserAgentSettings)
	cachedPrice             = make(map[string]*CachedPrice)
	cachedRSAPrivateKey     = make(map[string]*RSAPrivateKey)
	mx                      sync.Mutex
)

func InitUserAgentCache() {

	var userAgents []UserAgent

	DB.Set("gorm:auto_preload", true).Find(&userAgents)

	for _, userAgent := range userAgents {
		userAgent.CacheUserAgent()
	}
}

func (ua *UserAgent) BeforeCreate() (err error) {

	rsa_private_key, rsa_public_key := GenerateKeyPair(ServerConfig.RsaKeySize)

	rsa_private_key_bytes := PrivateKeyToBytes(rsa_private_key)
	rsa_public_key_bytes := PublicKeyToBytes(rsa_public_key)

	aes_encripted_rsa_private_key_bytes := AesEncrypt(rsa_private_key_bytes, ServerConfig.AesPassphrase)
	aes_encripted_rsa_public_key_bytes := AesEncrypt(rsa_public_key_bytes, ServerConfig.AesPassphrase)

	ua.RsaKey = RsaKey{
		RsaPrivateKeyAesEncripted: aes_encripted_rsa_private_key_bytes,
		RsaPublicKeyAesEncripted:  aes_encripted_rsa_public_key_bytes,
	}

	return
}

func (userAgent *UserAgent) Create() (err error) {
	if err = DB.Create(userAgent).Error; err != nil {
		return
	}
	DB.Save(userAgent)
	userAgent.CacheUserAgent()
	return
}

func (userAgent *UserAgent) Save() (err error) {
	DB.Save(userAgent)
	userAgent.CacheSettings()
	userAgent.CachePrice()
	return
}

func (userAgent *UserAgent) FindSettings() (err error) {

	if cUAS, ok := cachedUserAgentSettings[userAgent.Name]; !ok {
		if DB.First(userAgent).RecordNotFound() {
			err = ErrStatusForbidden
			return
		} else {
			userAgent.CacheSettings()
			return
		}
	} else {
		copier.Copy(userAgent, cUAS)
	}
	return
}

func (userAgent *UserAgent) FindPrice() (err error) {
	if cP, ok := cachedPrice[userAgent.Name]; !ok {
		if DB.First(userAgent).RecordNotFound() {
			err = ErrStatusForbidden
			return
		} else {
			userAgent.CachePrice()
			return
		}
	} else {
		copier.Copy(userAgent, cP)
	}
	return
}

func (userAgent *UserAgent) CacheSettings() (err error) {
	mx.Lock()
	defer mx.Unlock()
	cachedUserAgentSettings[userAgent.Name] = &CachedUserAgentSettings{
		Activitylimit: userAgent.Activitylimit,
		Like:          userAgent.Like,
		Follow:        userAgent.Follow,
		Pricefollow:   userAgent.Pricefollow,
		Pricelike:     userAgent.Pricelike,
	}
	return
}

func (userAgent *UserAgent) CachePrice() (err error) {
	mx.Lock()
	defer mx.Unlock()
	cachedPrice[userAgent.Name] = &CachedPrice{
		Pricefollow: userAgent.Pricefollow,
		Pricelike:   userAgent.Pricelike,
	}
	return
}

func (userAgent *UserAgent) CacheUserAgent() (err error) {
	userAgent.CacheSettings()
	userAgent.CachePrice()
	rsa_private_key_bytes := AesDecrypt(userAgent.RsaKey.RsaPrivateKeyAesEncripted, ServerConfig.AesPassphrase)
	cachedRSAPrivateKey[userAgent.Name] = BytesToPrivateKey(rsa_private_key_bytes)
	return
}
