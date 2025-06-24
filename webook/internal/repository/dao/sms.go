package dao

import (
	"context"
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
	FindByUnsend(ctx context.Context) (AsyncSms, error)
}
type GormAsyncSmsDAO struct {
	db *gorm.DB
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
func (d *GormAsyncSmsDAO) FindByUnsend(ctx context.Context) (AsyncSms, error) {
	var sms AsyncSms
	err := d.db.WithContext(ctx).Where("status=?", SmsWaiting).First(&sms).Error
	if err != nil {
		return AsyncSms{}, err
	}
	return sms, err
}

type AsyncSms struct {
	Id      int64 `gorm:"primaryKey,autoIncrement"`
	Biz     string
	Args    []string
	Numbers []string
	Status  int8
	Ctime   int64
	Utime   int64
}
