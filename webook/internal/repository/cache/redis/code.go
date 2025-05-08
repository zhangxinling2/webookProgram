package redis

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"webookProgram/webook/internal/repository/cache"
)

var (
	ErrSetCodeSendTooMany = errors.New("发送太频繁")
	ErrCodeVerifyTooMany  = errors.New("验证次数太多")
	ErrUnknown            = errors.New("未知错误")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeRedisCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) cache.CodeCache {
	return &CodeRedisCache{client: client}
}
func (c *CodeRedisCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		zap.L().Warn("短信发送太频繁", zap.String("biz", biz)) //要在对应的告警系统里面配置，比如说一分钟之内出现超过100次 WARN，就告警
		return ErrSetCodeSendTooMany
	default:
		return errors.New("系统错误")
	}
}
func (c *CodeRedisCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
func (c *CodeRedisCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return false, ErrUnknown
	}
}
