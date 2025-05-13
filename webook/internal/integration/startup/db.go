package startup

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webookProgram/webook/internal/repository/dao"
)

var db *gorm.DB

func InitTestDb() *gorm.DB {
	if db == nil {
		var err error
		db, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/webook"))
		if err != nil {
			panic(err)
		}
		err = dao.InitTable(db)
		if err != nil {
			panic(err)
		}
	}
	return db
}
