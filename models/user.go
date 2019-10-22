package models

import (
	"github.com/lib/pq"
	"time"
)

type User struct {
	Instagramid uint `json:"instagramid" binding:"required" gorm:"primary_key" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time     `sql:"index"`
	Banned      bool           `gorm:"default:false"`
	Coins       int            `json:"coins" gorm:"default:0"`
	Deviceid    string         `json:"deviceid" gorm:"-"`
	DeviceIds   pq.StringArray `gorm:"type:varchar(100)[]"`
	Rateus      bool           `json:"rateus" gorm:"default:true"`
}
