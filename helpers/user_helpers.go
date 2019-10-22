package helpers

import (
	"errors"
	"github.com/jinzhu/gorm"
	"instatasks/models"
)

var (
	ErrStatusForbidden = errors.New("Forbidden")
)

func FirstNotBannedUser(user *models.User, db *gorm.DB) error {

	tx := db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.First(&user).Error; err != nil {
		return err
	} else if user.Banned {
		return ErrStatusForbidden
	}
	if !tx.First(&models.BannedDevice{Deviceid: user.Deviceid}).RecordNotFound() {
		return ErrStatusForbidden
	}

	return tx.Commit().Error
}

func IsStatusForbiddenError(err error) bool {
	if errs, ok := err.(gorm.Errors); ok {
		for _, err := range errs {
			if err == ErrStatusForbidden {
				return true
			}
		}
	}
	return err == ErrStatusForbidden
}
