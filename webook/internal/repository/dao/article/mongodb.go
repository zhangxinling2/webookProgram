package article

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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
func (m MongoDBArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBArticleDAO) UpdateById(ctx context.Context, article Article) error {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBArticleDAO) Upsert(ctx context.Context, art PublishArticle) error {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBArticleDAO) Transaction(ctx context.Context, bizFunc func(txDAO ArticleDAO) error) error {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBArticleDAO) FindByID(ctx context.Context, articleId int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBArticleDAO) SyncStatus(ctx *gin.Context, id int64, author int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}
func InitCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

}
