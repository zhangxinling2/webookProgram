package article

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type AuthorDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}
type GORMAuthorDAO struct {
	db *gorm.DB
}

func (a *GORMAuthorDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.CTime = now
	art.CTime = now
	err := a.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (a *GORMAuthorDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.UTime = now
	res := a.db.WithContext(ctx).Model(&art).Where("id=?", art.Id).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"u_time":  art.UTime,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func NewGORMAuthorDAO(db *gorm.DB) AuthorDAO {
	return &GORMAuthorDAO{
		db: db,
	}
}
