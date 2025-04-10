package repository

import (
	"context"
	"database/sql"
	"time"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/cache"
	"webookProgram/webook/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByWechat(ctx context.Context, openId string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	UpdateInfo(ctx context.Context, u domain.User) (domain.User, error)
}
type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}
func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}
func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}
func (r *CacheUserRepository) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	u, err := r.dao.FindByOpenId(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}
func (r *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}
	user, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.entityToDomain(user)
	err = r.cache.Set(ctx, u)
	if err != nil {
		//打个日志,做监控即可
	}
	return u, err
}
func (r *CacheUserRepository) UpdateInfo(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := r.dao.EditUser(ctx, r.domainToEntity(u))
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(user), err
}
func (r *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:           u.Id,
		Email:        u.Email.String,
		Password:     u.Password,
		Phone:        u.Phone.String,
		NickName:     u.NickName,
		Birth:        time.UnixMilli(u.Birth),
		Introduction: u.Introduction,
		CTime:        time.UnixMilli(u.Ctime),
		WechatInfo: domain.WechatInfo{
			UnionID: u.WechatUnionId.String,
			OpenID:  u.WechatOpenId.String,
		},
	}
}
func (r *CacheUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		WechatOpenId: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionId: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
		Password:     u.Password,
		NickName:     u.NickName,
		Birth:        u.Birth.UnixMilli(),
		Introduction: u.Introduction,
		Ctime:        u.CTime.UnixMilli(),
	}
}
