package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository"
	svcmock "webookProgram/webook/internal/repository/mocks"
)

func Test_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name:     "登录成功",
			email:    "123@qq.com",
			password: "hello#123",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := svcmock.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$xOJ28aKxNGcBqw0vXE7yUepBkkCY1FcFL1dXkvtKBX678Eui4vl6G",
						Phone:    "11111111111",
						CTime:    now,
					}, nil)
				return repo
			},
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$xOJ28aKxNGcBqw0vXE7yUepBkkCY1FcFL1dXkvtKBX678Eui4vl6G",
				Phone:    "11111111111",
				CTime:    now,
			},
		},
		{
			name:     "用户不存在",
			email:    "123@qq.com",
			password: "hello#123",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := svcmock.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name:     "DB错误",
			email:    "123@qq.com",
			password: "hello#123",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := svcmock.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name:     "密码不对",
			email:    "123@qq.com",
			password: "123hello#123",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := svcmock.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$xOJ28aKxNGcBqw0vXE7yUepBkkCY1FcFL1dXkvtKBX678Eui4vl6G",
						Phone:    "11111111111",
						CTime:    now,
					}, nil)
				return repo
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			usr, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, usr)
		})
	}
}
func TestEncrypted(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("hello#123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(hash))
	}
}
