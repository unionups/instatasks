package models

import "time"

type BannedDevice struct {
	Deviceid  string `gorm:"primary_key:true"`
	CreatedAt time.Time
}
