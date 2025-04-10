package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账户/邮箱或密码不对")
)

type UserService interface {
	FindOrCreate(ctx *gin.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx *gin.Context, wechat domain.WechatInfo) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
func (svc *userService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	//要判断有没有这个用户
	//这里是快路径
	if err != repository.ErrUserNotFound {
		return u, err
	}
	//在系统资源不足，触发降级之后，不执行慢路径了
	//这里是慢路径
	err = svc.repo.Create(ctx, domain.User{Phone: phone})
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	//存在主从延迟
	return svc.repo.FindByPhone(ctx, phone)
}
func (svc *userService) FindOrCreateByWechat(ctx *gin.Context, wechat domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechat.OpenID)

	if err != repository.ErrUserNotFound {
		return u, err
	}

	err = svc.repo.Create(ctx, domain.User{
		WechatInfo: wechat,
	})
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	//存在主从延迟
	return svc.repo.FindByWechat(ctx, wechat.OpenID)
}
func (svc *userService) Edit(ctx context.Context, u domain.User) error {
	_, err := svc.repo.UpdateInfo(ctx, u)
	if err != nil {
		return err
	}
	return err
}
func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {

	u, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}
func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
