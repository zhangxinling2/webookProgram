package startup

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"webookProgram/webook/internal/repository/dao"
)

var db *gorm.DB
var mongoDb *mongo.Database

func InitNode() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}
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
func InitMongoDB() *mongo.Database {
	if mongoDb == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context,
				startedEvent *event.CommandStartedEvent) {
				fmt.Println(startedEvent.Command)
			},
		}
		opt := options.Client().ApplyURI("mongodb://localhost:27017").
			SetTimeout(time.Second * 10).
			SetAuth(options.Credential{Username: "root", Password: "123456"}).SetMonitor(monitor)
		client, err := mongo.Connect(ctx, opt)
		if err != nil {
			panic(err)
		}
		return client.Database("test")
	}
	return mongoDb
}
