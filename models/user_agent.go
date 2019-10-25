package models

import (
	"time"
)

type UserAgent struct {
	Name      string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Activitylimit uint `json:"activitylimit" gorm:"default:0"`
	Like          bool `json:"like" gorm:"default:true"`
	Follow        bool `json:"follow" gorm:"default:true"`
	Pricefollow   uint `json:"pricefollow" gorm:"default:5"`
	Pricelike     uint `json:"pricefollow" gorm:"default:1"`

	RsaPrivateKeyAesEncripted []byte `gorm:"type:byte[];not null;"`
	RsaPublicKeyAesEncripted  []byte `gorm:"type:byte[];not null;"`
}
