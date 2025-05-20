package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	FindByID(ctx context.Context, id int64) (Article, error)
	Update(ctx context.Context, art Article) error
}
type GORMArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}
func (d *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.CTime = now
	art.UTime = now
	err := d.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}
func (d *GORMArticleDAO) Update(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.UTime = now
	res := d.db.Model(&art).WithContext(ctx).Where("id=? and author_id=?", art.Id, art.AuthorId).Updates(map[string]any{
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
func (d *GORMArticleDAO) FindByID(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := d.db.WithContext(ctx).Where("id=?", id).First(&art).Error
	if err != nil {
		return Article{}, err
	}
	return art, nil
}

type Article struct {
	Id       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Title    string `gorm:"type:varchar(1024)"`
	Content  string `gorm:"type:Blob"`
	AuthorId int64  `gorm:"index"`
	CTime    int64
	UTime    int64
}
