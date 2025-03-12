package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/cache"
	cachemock "webookProgram/webook/internal/repository/cache/mocks"
	"webookProgram/webook/internal/repository/dao"
	daomock "webookProgram/webook/internal/repository/dao/mocks"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	//要去掉毫秒以外的部分
	//now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemock.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrUserNotFound)
				d := daomock.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{
					Id:       1,
					Password: "xxx",
					Phone: sql.NullString{
						String: "123123123",
						Valid:  true,
					},
					Birth: now.UnixMilli(),
					Email: sql.NullString{
						Valid:  true,
						String: "123@qq.com",
					},
					Ctime: now.UnixMilli(),
					Utime: now.UnixMilli(),
				}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       1,
					Password: "xxx",
					Phone:    "123123123",
					Birth:    now,
					Email:    "123@qq.com",
					CTime:    now,
				}).Return(nil)
				return d, c
			},
			id: 1,
			wantUser: domain.User{
				Id:       1,
				Password: "xxx",
				Birth:    now,
				Phone:    "123123123",
				Email:    "123@qq.com",
				CTime:    now,
			},
		},
		{
			name: "缓存命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemock.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{
					Id:       1,
					Password: "xxx",
					Birth:    now,
					Phone:    "123123123",
					Email:    "123@qq.com",
					CTime:    now,
				}, nil)
				d := daomock.NewMockUserDAO(ctrl)
				return d, c
			},
			id: 1,
			wantUser: domain.User{
				Id:       1,
				Password: "xxx",
				Birth:    now,
				Phone:    "123123123",
				Email:    "123@qq.com",
				CTime:    now,
			},
		},
		{
			name: "缓存未命中，查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemock.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrUserNotFound)
				d := daomock.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{}, errors.New("mock db 错误"))

				return d, c
			},
			id:       1,
			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}
