package repository

import (
	"context"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/dao"
)

var (
	ErrNoSms = dao.ErrNoSms
)

type AsyncSmsRepository interface {
	Store(ctx context.Context, sms domain.AsyncSms) error
	PrepareWaitingSms(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, id int64, status bool) error
}
type asyncSmsRepository struct {
	dao dao.AsyncSmsDAO
}

func (r *asyncSmsRepository) ReportScheduleResult(ctx context.Context, id int64, status bool) error {
	if status {
		return r.dao.MarkSuccess(ctx, id)
	}
	return r.dao.MarkFail(ctx, id)
}

func NewAsyncSmsReposity(dao dao.AsyncSmsDAO) AsyncSmsRepository {
	return &asyncSmsRepository{
		dao: dao,
	}
}
func (r *asyncSmsRepository) Store(ctx context.Context, sms domain.AsyncSms) error {
	return r.dao.Insert(ctx, r.domainToDao(sms))
}
func (r *asyncSmsRepository) PrepareWaitingSms(ctx context.Context) (domain.AsyncSms, error) {
	smsInfo, err := r.dao.GetWaitingSms(ctx)
	if err != nil {
		return domain.AsyncSms{}, err
	}
	return r.daoToDomain(smsInfo), nil
}
func (r *asyncSmsRepository) domainToDao(sms domain.AsyncSms) dao.AsyncSms {
	return dao.AsyncSms{
		Biz:      sms.Biz,
		Args:     sms.Args,
		Numbers:  sms.Numbers,
		RetryMax: sms.RetryMax,
	}
}
func (r *asyncSmsRepository) daoToDomain(sms dao.AsyncSms) domain.AsyncSms {
	return domain.AsyncSms{
		Id:      sms.Id,
		Biz:     sms.Biz,
		Args:    sms.Args,
		Numbers: sms.Numbers,
	}
}
