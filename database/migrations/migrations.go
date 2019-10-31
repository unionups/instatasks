package migrations

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/gormigrate.v1"
	"instatasks/database"
	"instatasks/models"

	"log"
	"time"
)

func Migrate() {
	db := database.GetDB()
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create users table
		{
			ID: "101608301601",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time
				type User struct {
					Instagramid uint `json:"instagramid" binding:"required" gorm:"primary_key:true"`
					CreatedAt   time.Time
					UpdatedAt   time.Time
					DeletedAt   *time.Time `sql:"index"`
					Banned      bool       `gorm:"default:false"`
					Coins       uint       `json:"coins" gorm:"default:0"`
					Deviceid    string     `json:"deviceid" gorm:"-"`
					Rateus      bool       `binding:"-" gorm:"default:true"`

					Tasks []models.Task `gorm:"foreignkey:Instagramid;association_foreignkey:Instagramid;"`
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
					Deviceid  string `gorm:"primary_key:true"`
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
			ID: "101608301802",
			Migrate: func(tx *gorm.DB) error {
				type UserAgent struct {
					Name      string `header:"User-Agent" json:"name" binding:"required"  gorm:"primary_key:true"`
					CreatedAt time.Time
					UpdatedAt time.Time
					DeletedAt *time.Time

					Activitylimit uint `json:"activitylimit" gorm:"default:0"`
					Like          bool `json:"like" gorm:"default:true"`
					Follow        bool `json:"follow" gorm:"default:true"`
					Pricefollow   uint `json:"pricefollow" gorm:"default:5"`
					Pricelike     uint `json:"pricelike" gorm:"default:1"`

					RsaKey models.RsaKey `gorm:"foreignkey:Name;association_foreignkey:Name"`
				}
				return tx.AutoMigrate(&UserAgent{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("user_agents").Error
			},
		},
		// create rsa_keys table
		{
			ID: "101608301901",
			Migrate: func(tx *gorm.DB) error {
				type RsaKey struct {
					Name string `header:"User-Agent" json:"name" gorm:"primary_key:true"`

					RsaPrivateKeyAesEncripted []byte `gorm:"not null"`
					RsaPublicKeyAesEncripted  []byte `gorm:"not null"`
				}
				return tx.AutoMigrate(&RsaKey{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("rsa_keys").Error
			},
		},
		// create tasks table
		{
			ID: "101608302001",
			Migrate: func(tx *gorm.DB) error {

				type Task struct {
					ID                uint       `json:"-" gorm:"primary_key"`
					CreatedAt         time.Time  `json:"created_at"`
					DeletedAt         *time.Time `json:"deleted_at" sql:"index"`
					Taskid            string     `json:"taskid" gorm:"-"`
					Type              string     `json:"type" binding:"required"`
					Count             uint       `json:"count" binding:"required" gorm:"not null"`
					LeftCounter       uint       `json:"left_counter"`
					Photourl          string     `json:"photourl"`
					Instagramusername string     `json:"instagramusername"`
					Mediaid           string     `json:"mediaid" binding:"required" sql:"index" gorm:"not null"`

					CancelLeftCounter uint8 `json:"-"`

					Instagramid uint `json:"instagramid" binding:"required" sql:"index" gorm:"not null"`
				}

				return tx.AutoMigrate(&Task{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("tasks").Error
			},
		},
		// create user_mediaids table
		{
			ID: "101608302101",
			Migrate: func(tx *gorm.DB) error {
				type UserMediaid struct {
					Instagramid uint   `sgl:"index"`
					Mediaid     string `sql:"index"`
				}
				return tx.AutoMigrate(&UserMediaid{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("user_mediaids").Error
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")

}
