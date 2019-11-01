package models

type UserMediaid struct {
	ID          uint
	Instagramid uint   `sql:"index"`
	Mediaid     string `sql:"index"`
}
