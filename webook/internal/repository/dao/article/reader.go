package article

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReaderDAO interface {
	UpsertArticle(ctx context.Context, art Article) error
}
type ReaderGORMDAO struct {
	db *gorm.DB
}

func (r *ReaderGORMDAO) UpsertArticle(ctx context.Context, art Article) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
		}),
	}).Error
}

func NewReaderGORMDAO(db *gorm.DB) ReaderDAO {
	return &ReaderGORMDAO{db: db}
}

type PublishArticle struct {
	Article
}
