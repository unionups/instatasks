package helpers

import (
	"errors"
	"github.com/jinzhu/gorm"
)

var (
	ErrStatusForbidden = errors.New("Forbidden")
)

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
