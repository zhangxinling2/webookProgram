package cache

import (
	"context"
	_ "embed"
	"errors"
)

var (
	ErrSetCodeSendTooMany = errors.New("发送太频繁")
	ErrCodeVerifyTooMany  = errors.New("验证次数太多")
	ErrUnknown            = errors.New("未知错误")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}
