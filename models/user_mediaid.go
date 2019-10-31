package models

type UserMediaid struct {
	Instagramid uint   `sgl:"index"`
	Mediaid     string `sql:"index"`
}
