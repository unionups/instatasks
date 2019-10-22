package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"instatasks/config"
	"os"
	"time"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB
var err error

func InitDB() *gorm.DB {
	var db = DB
	var dbURL string

	config := config.InitConfig()

	fmt.Println(config.AppEnv)
	if config.AppEnv == "staging" {
		dbURL = os.Getenv("DATABASE_URL")
	} else {
		// driver := config.Database.Driver
		database := config.Database.Dbname
		username := config.Database.Username
		password := config.Database.Password
		host := config.Database.Host
		port := config.Database.Port
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			username,
			password,
			host,
			port,
			database,
		)
	}

	db, err = gorm.Open("postgres", dbURL)

	if err != nil {
		panic(err)
		panic("failed to connect database")
	}

	db.DB().SetMaxIdleConns(config.Database.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Hour)

	if config.AppEnv == "development" {
		db.LogMode(true)
	}

	DB = db

	return DB
}

func GetDB() *gorm.DB {
	return DB
}

func Migrate() {
	migrate()
}
