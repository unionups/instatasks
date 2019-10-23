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
					Instagramid uint `json:"instagramid" binding:"required" gorm:"primary_key" `
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
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")

}
