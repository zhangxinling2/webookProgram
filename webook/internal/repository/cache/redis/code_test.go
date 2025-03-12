package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	redismock "webookProgram/webook/internal/repository/cache/redismocks"
)

func TestCodeRedisCache_Set(t *testing.T) {
	testCases := []struct {
		name  string
		mock  func(ctrl *gomock.Controller) redis.Cmdable
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).
					Return(res)

				return cmd
			},
			biz:   "login",
			phone: "152",
			code:  "123456",
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock redis 错误"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).
					Return(res)

				return cmd
			},
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: errors.New("mock redis 错误"),
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).
					Return(res)

				return cmd
			},
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: ErrSetCodeSendTooMany,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-2))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).
					Return(res)

				return cmd
			},
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: errors.New("系统错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
