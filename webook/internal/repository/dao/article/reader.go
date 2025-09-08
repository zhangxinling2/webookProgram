package article

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReaderDAO interface {
	//操作不同库
	UpsertArticle(ctx context.Context, art Article) error
	//同库不同表
	UpsertArticleV1(ctx context.Context, art PublishedArticle) error
}
type ReaderGORMDAO struct {
	db *gorm.DB
}

// 不同库
func (r *ReaderGORMDAO) UpsertArticle(ctx context.Context, art Article) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
		}),
	}).Create(&art).Error
}

// 同库不同表
func (r *ReaderGORMDAO) UpsertArticleV1(ctx context.Context, art PublishedArticle) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
		}),
	}).Create(&art).Error
}
func NewReaderGORMDAO(db *gorm.DB) ReaderDAO {
	return &ReaderGORMDAO{db: db}
}

type PublishedArticle struct {
	Article
}
