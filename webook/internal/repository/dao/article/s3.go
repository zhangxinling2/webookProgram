package article

import (
	"bytes"
	"context"
	_ "github.com/aws/aws-sdk-go-v2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
	"strconv"
	"webookProgram/ekit"
	"webookProgram/webook/internal/domain"
)

type S3DAO struct {
	oss *s3.Client
	GORMArticleDAO
	bucket *string
}

func NewOssDAO(oss *s3.Client, db *gorm.DB) *S3DAO {
	return &S3DAO{oss: oss,
		GORMArticleDAO: GORMArticleDAO{db},
		bucket:         ekit.ToPtr[string]("webook-1314583317"),
	}
}
func (o *S3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id = art.Id
	)
	err := o.db.Transaction(func(tx *gorm.DB) error {
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
		art.Id = id
		art.Content = ""
		// 操作线上库了
		return txDAO.Upsert(ctx, PublishedArticle{Article: art})
	})
	if err != nil {
		return 0, err
	}
	_, err = o.oss.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      o.bucket,
		Key:         ekit.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	return art.Id, err
}
func (o *S3DAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id=? AND author_id = ?", id, author).
			Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}

		res = tx.Model(&PublishedArticle{}).
			Where("id=? AND author_id = ?", id, author).Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}
		return nil
	})
	if err != nil {
		return err
	}
	if status == domain.ArticlePrivate.ToUint8() {
		_, err = o.oss.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: o.bucket,
			Key:    ekit.ToPtr[string](strconv.FormatInt(id, 10)),
		})
	}
	return err
}
