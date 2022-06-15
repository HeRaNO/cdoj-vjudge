package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var RDB *gorm.DB

func InitDB() {
	var err error
	config := conf.RDB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		config.Username, config.Password, config.Host, config.Port, config.DBName, config.TimeZone)
	RDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("[FAILED] init RDB failed")
		panic(err)
	}
	log.Println("[INFO] init database finished successfully")
}
