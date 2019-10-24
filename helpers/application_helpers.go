package helpers

import (
	"reflect"
	"time"
)

func IsItemInSlice(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("isItemInSlice() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func DurationInHours(i int) time.Duration {
	return time.Duration(i) * time.Hour
}
