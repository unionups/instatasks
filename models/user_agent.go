package models

import (
	"instatasks/config"
	. "instatasks/helpers"
	"time"
)

type UserAgent struct {
	Name      string `header:"User-Agent" json:"name" binding:"required"  gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Activitylimit uint `json:"activitylimit" gorm:"default:0"`
	Like          bool `json:"like" gorm:"default:true"`
	Follow        bool `json:"follow" gorm:"default:true"`
	Pricefollow   uint `json:"pricefollow" gorm:"default:5"`
	Pricelike     uint `json:"pricelike" gorm:"default:1"`

	RsaKey RsaKey `gorm:"foreignkey:Name;association_foreignkey:Name"`
}

func (ua *UserAgent) BeforeCreate() (err error) {

	serverConfig := config.GetConfig().Server

	rsa_private_key, rsa_public_key := GenerateKeyPair(serverConfig.RsaKeySize)

	rsa_private_key_bytes := PrivateKeyToBytes(rsa_private_key)
	rsa_public_key_bytes := PublicKeyToBytes(rsa_public_key)

	aes_encripted_rsa_private_key_bytes := AesEncrypt(rsa_private_key_bytes, serverConfig.AesPassphrase)
	aes_encripted_rsa_public_key_bytes := AesEncrypt(rsa_public_key_bytes, serverConfig.AesPassphrase)

	ua.RsaKey = RsaKey{
		RsaPrivateKeyAesEncripted: aes_encripted_rsa_private_key_bytes,
		RsaPublicKeyAesEncripted:  aes_encripted_rsa_public_key_bytes,
	}

	return
}

// func (ua *UserAgent) BeforeCreate() (err error) {

// 	serverConfig := config.GetConfig().Server

// 	rsa_private_key, rsa_public_key := GenerateKeyPair(serverConfig.RsaKeySize)

// 	rsa_private_key_bytes := PrivateKeyToBytes(rsa_private_key)
// 	rsa_public_key_bytes := PublicKeyToBytes(rsa_public_key)

// 	aes_encripted_rsa_private_key_bytes := AesEncrypt(rsa_private_key_bytes, serverConfig.AesPassphrase)
// 	aes_encripted_rsa_public_key_bytes := AesEncrypt(rsa_public_key_bytes, serverConfig.AesPassphrase)

// 	ua.RsaKey = RsaKey{
// 		RsaPrivateKeyAesEncripted: aes_encripted_rsa_private_key_bytes,
// 		RsaPublicKeyAesEncripted:  aes_encripted_rsa_public_key_bytes,
// 	}

// 	return
// }
