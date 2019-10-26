package models

type RsaKey struct {
	Name string `header:"User-Agent" json:"name"  gorm:"primary_key"`

	RsaPrivateKeyAesEncripted []byte `gorm:"type:byte[];not null;"`
	RsaPublicKeyAesEncripted  []byte `gorm:"type:byte[];not null;"`
}
