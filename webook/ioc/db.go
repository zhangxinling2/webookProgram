package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webookProgram/webook/config"
	"webookProgram/webook/internal/repository/dao"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
