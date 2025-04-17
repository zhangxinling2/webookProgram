package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webookProgram/webook/internal/repository/dao"
)

func InitDb() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var config Config
	err := viper.UnmarshalKey("db", &config)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(config.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
