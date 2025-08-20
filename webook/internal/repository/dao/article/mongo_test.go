package article

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opt := options.Client().ApplyURI("mongodb://localhost:27017").
		SetTimeout(time.Second * 10).
		SetAuth(options.Credential{Username: "root", Password: "123456"})
	client, err := mongo.Connect(ctx, opt)

	assert.NoError(t, err)
	defer func() {
		// 确保在函数结束时断开与MongoDB的连接
		if err := client.Disconnect(ctx); err != nil {
			panic(err) // 如果断开连接失败，立即终止程序
		}
	}()
	coll := client.Database("test").Collection("test")
	res, err := coll.InsertOne(ctx, Article{
		Id:       123,
		Title:    "test",
		Content:  "test",
		AuthorId: 111,
		Status:   1,
		CTime:    123,
		UTime:    123,
	})
	assert.NoError(t, err)
	fmt.Println(res.InsertedID)
}

type MongoArticle struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	CTime    int64  `bson:"ctime,omitempty"`
	UTime    int64  `bson:"utime,omitempty"`
}
