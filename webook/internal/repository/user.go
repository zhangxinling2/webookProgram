package repository

import (
	"context"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	})
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
}
func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}
	user, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:           user.Id,
		Email:        user.Email,
		Password:     user.Password,
		NickName:     user.NickName,
		Birth:        time.UnixMilli(user.Birth),
		Introduction: user.Introduction,
		CTime:        time.UnixMilli(user.Ctime),
	}
	err = r.cache.Set(ctx, u)
	if err != nil {
		//打个日志,做监控即可
	}
	return u, err
}
func (r *UserRepository) UpdateInfo(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := r.dao.EditUser(ctx, dao.User{
		Id:           u.Id,
		Email:        u.Email,
		NickName:     u.NickName,
		Birth:        u.Birth.UnixMilli(),
		Introduction: u.Introduction,
		Password:     u.Password,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:           user.Id,
		Email:        user.Email,
		Password:     user.Password,
		NickName:     user.NickName,
		Birth:        time.UnixMilli(user.Birth),
		Introduction: user.Introduction,
		CTime:        time.UnixMilli(user.Ctime),
	}, err
}
