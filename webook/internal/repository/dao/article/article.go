package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticle) error
	Transaction(ctx context.Context,
		bizFunc func(txDAO ArticleDAO) error) error
	FindByID(ctx context.Context, articleId int64) (Article, error)
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleDAO) FindByID(ctx context.Context, articleId int64) (Article, error) {
	var art Article
	err := dao.db.WithContext(ctx).Where("id=?", articleId).First(&art).Error
	return art, err
}
func (dao *GORMArticleDAO) Transaction(ctx context.Context,
	bizFunc func(txDAO ArticleDAO) error) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMArticleDAO(tx)
		return bizFunc(txDAO)
	})
}

// 同库不同表
func (dao *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 先操作制作库（此时应该是表），后操作线上库（此时应该是表）

	var (
		id = art.Id
	)
	// tx => Transaction, trx, txn
	// 在事务内部，这里采用了闭包形态
	// GORM 帮助我们管理了事务的生命周期
	// Begin，Rollback 和 Commit 都不需要我们操心
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewGORMArticleDAO(tx)
		if id > 0 {
			err = txDAO.UpdateById(ctx, art)
		} else {
			id, err = txDAO.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		// 操作线上库了
		return txDAO.Upsert(ctx, PublishArticle{Article: art})
	})
	return id, err
}

// Upsert INSERT OR UPDATE
func (dao *GORMArticleDAO) Upsert(ctx context.Context, art PublishArticle) error {
	now := time.Now().UnixMilli()
	art.CTime = now
	art.UTime = now
	// 这个是插入
	// OnConflict 的意思是数据冲突了
	err := dao.db.Clauses(clause.OnConflict{
		// SQL 2003 标准
		// INSERT AAAA ON CONFLICT(BBB) DO NOTHING
		// INSERT AAAA ON CONFLICT(BBB) DO UPDATES CCC WHERE DDD

		// 哪些列冲突
		//Columns: []clause.Column{clause.Column{Name: "id"}},
		// 意思是数据冲突，啥也不干
		// DoNothing:
		// 数据冲突了，并且符合 WHERE 条件的就会执行 DO UPDATES
		// Where:

		// MySQL 只需要关心这里
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"u_time":  now,
		}),
	}).Create(&art).Error
	// MySQL 最终的语句 INSERT xxx ON DUPLICATE KEY UPDATE xxx

	// 一条 SQL 语句，都不需要开启事务
	// auto commit: 意思是自动提交

	return err
}

// 事务传播机制是指如果当前有事务，就在事务内部执行 Insert
// 如果咩有事务：
// 1. 开启事务，执行 Insert
// 2. 直接执行
// 3. 报错

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.CTime = now
	art.UTime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.UTime = now
	// 依赖 gorm 忽略零值的特性，会用主键进行更新
	// 可读性很差
	res := dao.db.WithContext(ctx).Model(&art).
		Where("id=? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"u_time":  art.UTime,
		})
	// 你要不要检查真的更新了没？
	// res.RowsAffected // 更新行数
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		//dangerousDBOp.Count(1)
		// 补充一点日志
		return fmt.Errorf("更新失败，可能是创作者非法 id %d, author_id %d",
			art.Id, art.AuthorId)
	}
	return res.Error
}

type Article struct {
	Id       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Title    string `gorm:"type:varchar(1024)"`
	Content  string `gorm:"type:Blob"`
	AuthorId int64  `gorm:"index"`
	CTime    int64
	UTime    int64
}
