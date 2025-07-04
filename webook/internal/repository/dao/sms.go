package dao

import (
	"context"
	"gorm.io/gorm/clause"
	"time"

	"gorm.io/gorm"
)

const (
	SmsWaiting = iota
	SmsFailed
	SmsSuccess
)

var (
	ErrNoSms = gorm.ErrRecordNotFound
)

type AsyncSmsDAO interface {
	Insert(ctx context.Context, sms AsyncSms) error
	GetWaitingSms(ctx context.Context) (AsyncSms, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFail(ctx context.Context, id int64) error
}
type GormAsyncSmsDAO struct {
	db *gorm.DB
}

func (d *GormAsyncSmsDAO) MarkSuccess(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return d.db.WithContext(ctx).Model(&AsyncSms{}).Where("id=?", id).Updates(map[string]any{
		"utime":  now,
		"status": SmsSuccess,
	}).Error
}

func (d *GormAsyncSmsDAO) MarkFail(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return d.db.WithContext(ctx).Model(&AsyncSms{}).Where("id=?", id).Updates(map[string]any{
		"utime":  now,
		"status": SmsFailed,
	}).Error
}

func NewGormAsyncSmsDAO(db *gorm.DB) AsyncSmsDAO {
	return &GormAsyncSmsDAO{
		db: db,
	}
}
func (d *GormAsyncSmsDAO) Insert(ctx context.Context, sms AsyncSms) error {
	now := time.Now().UnixMilli()
	sms.Utime = now
	sms.Ctime = now
	return d.db.WithContext(ctx).Create(&sms).Error
}
func (d *GormAsyncSmsDAO) GetWaitingSms(ctx context.Context) (AsyncSms, error) {
	var sms AsyncSms
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		endtime := now - time.Minute.Milliseconds()
		//查询并锁住记录，只查询超过一分钟的，相当于一分钟重试
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("utime < ? and status = ?", endtime, SmsWaiting).First(&sms).Error
		if err != nil {
			return err
		}
		//更新时间防止其他服务调用
		err = tx.Model(&AsyncSms{}).Where("id = ?", sms.Id).Updates(map[string]any{
			"utime":     now,
			"retry_cnt": gorm.Expr("retry_cnt +1"),
		}).Error
		return err
	})
	return sms, err
}

type AsyncSms struct {
	Id       int64 `gorm:"primaryKey,autoIncrement"`
	Biz      string
	Args     []string
	Numbers  []string
	RetryCnt int8
	RetryMax int8
	Status   int8
	Ctime    int64
	Utime    int64
}
