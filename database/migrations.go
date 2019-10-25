package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"gopkg.in/gormigrate.v1"
	"time"
	// "instatasks/models"
)

func migrate() {
	db := DB
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create users table
		{
			ID: "101608301601",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time
				type User struct {
					Instagramid uint64 `json:"instagramid" binding:"required" gorm:"primary_key" `
					CreatedAt   time.Time
					UpdatedAt   time.Time
					DeletedAt   *time.Time `sql:"index"`
					Banned      bool       `gorm:"default:false"`
					Coins       int        `json:"coins" gorm:"default:0"`
					Deviceid    string     `json:"deviceid" gorm:"-"`
					Rateus      bool       `json:"rateus" gorm:"default:true"`
				}
				return tx.AutoMigrate(&User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("users").Error
			},
		},
		// create banned_devices table
		{
			ID: "101608301701",
			Migrate: func(tx *gorm.DB) error {
				type BannedDevice struct {
					Deviceid  string `gorm:"primary_key"`
					CreatedAt time.Time
				}
				return tx.AutoMigrate(&BannedDevice{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("banned_devices").Error
			},
		},
		// create user_agents table
		{
			ID: "101608301801",
			Migrate: func(tx *gorm.DB) error {
				type UserAgent struct {
					Name      string `gorm:"primary_key"`
					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt *time.Time `sql:"index"`

					Activitylimit uint `json:"activitylimit" gorm:"default:0"`
					Like          bool `json:"like" gorm:"default:true"`
					Follow        bool `json:"follow" gorm:"default:true"`
					Pricefollow   uint `json:"pricefollow" gorm:"default:5"`
					Pricelike     uint `json:"pricefollow" gorm:"default:1"`

					RsaPrivateKeyAesEncripted []byte `gorm:"type:byte[];not null;"`
					RsaPublicKeyAesEncripted  []byte `gorm:"type:byte[];not null;"`
				}
				return tx.AutoMigrate(&UserAgent{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("user_agents").Error
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")

}
