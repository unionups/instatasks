package models

import (
	"github.com/jinzhu/gorm"
	"instatasks/config"
	"instatasks/database"
)

var ServerConfig *config.ServerConfiguration
var DB *gorm.DB

func Init() {
	ServerConfig = &config.GetConfig().Server
	DB = database.GetDB()
	InitUserCache()
	InitUserAgentCache()
	InitTaskCache()
}
