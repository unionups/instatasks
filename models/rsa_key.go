package models

type RsaKey struct {
	Name string `header:"User-Agent" json:"name" gorm:"primary_key:true"`

	RsaPrivateKeyAesEncripted []byte `gorm:"not null"`
	RsaPublicKeyAesEncripted  []byte `gorm:"not null"`
}
