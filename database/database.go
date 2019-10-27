package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"instatasks/config"
	"os"
	"time"
)

var DB *gorm.DB
var err error

func InitDB() *gorm.DB {
	var db = DB
	var dbURL string

	config := config.InitConfig()
	appEnv := config.AppEnv

	fmt.Println(config.AppEnv)
	if appEnv == "staging" {
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

	if appEnv == "test" {
		db, err = gorm.Open("sqlite3", "file::memory:?cache=shared")
	} else {
		db, err = gorm.Open("postgres", dbURL)
	}

	if err != nil {
		panic(err)
		panic("failed to connect database")
	}

	if appEnv != "test" {
		db.DB().SetMaxIdleConns(config.Database.MaxIdleConns)
		db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
		db.DB().SetConnMaxLifetime(time.Hour)
	}

	if appEnv == "development" || appEnv == "test" {
		db.LogMode(true)
	}

	DB = db
	migrate()

	return DB
}

func GetDB() *gorm.DB {
	return DB
}
