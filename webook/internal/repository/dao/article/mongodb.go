package article

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBArticleDAO struct {
	col     *mongo.Collection
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func NewMongoDBArticleDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBArticleDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node:    node,
	}
}
func (m *MongoDBArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	art.Id = m.node.Generate().Int64()
	now := time.Now().UnixMilli()
	art.CTime = now
	art.UTime = now
	_, err := m.col.InsertOne(ctx, art)
	return art.Id, err
}

func (m *MongoDBArticleDAO) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.UTime = now
	filter := bson.M{"id": article.Id, "author_id": article.AuthorId}
	sets := bson.M{"$set": bson.M{"title": article.Title,
		"content": article.Content,
		"u_time":  article.UTime},
	}
	res, err := m.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.MatchedCount != 1 {
		// 比较可能就是有人更新别人的文章，比如说攻击者跟你过不去
		return errors.New("更新失败")
	}
	return nil
}

func (m *MongoDBArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, article)
	} else {
		id, err = m.Insert(ctx, article)
	}

}

func (m *MongoDBArticleDAO) Upsert(ctx context.Context, art PublishedArticle) error {
	//now := time.Now().UnixMilli()
	//art.UTime = now
	//filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	//_, err := m.liveCol.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
	//	"title":   art.Title,
	//	"content": art.Content,
	//	"u_time":  art.UTime},
	//	"$setOnInsert":})
}

func (m *MongoDBArticleDAO) Transaction(ctx context.Context, bizFunc func(txDAO ArticleDAO) error) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) FindByID(ctx context.Context, articleId int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) SyncStatus(ctx *gin.Context, id int64, author int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}
func InitCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "c_time", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").Indexes().
		CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_articles").Indexes().
		CreateMany(ctx, index)
	return err
}
