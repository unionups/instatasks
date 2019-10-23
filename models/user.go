package models

import (
	"github.com/jinzhu/gorm"
	. "instatasks/helpers"
	"time"
)

type User struct {
	Instagramid uint `json:"instagramid" binding:"required" gorm:"primary_key" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Banned      bool       `gorm:"default:false"`
	Coins       int        `json:"coins" gorm:"default:0"`
	Deviceid    string     `json:"deviceid" gorm:"-"`
	Rateus      bool       `json:"rateus" gorm:"default:true"`
}

func FirstNotBannedUserScope(user *User, db *gorm.DB) error {

	tx := db.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	if err := tx.First(&user).Error; err != nil {
		if !tx.First(&BannedDevice{Deviceid: user.Deviceid}).RecordNotFound() {
			tx.Rollback()
			return ErrStatusForbidden
		}
		tx.Rollback()
		return err
	} else if user.Banned {
		tx.Rollback()
		return ErrStatusForbidden
	}

	if !tx.First(&BannedDevice{Deviceid: user.Deviceid}).RecordNotFound() {
		tx.Rollback()
		return ErrStatusForbidden
	}

	return tx.Commit().Error
}
