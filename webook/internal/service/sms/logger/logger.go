package logger

import (
	"context"
	"go.uber.org/zap"
	"webookProgram/webook/internal/service/sms"
)

type Service struct {
	svc sms.Service
}

func (s Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	zap.L().Debug("发送短信", zap.String("biz", biz), zap.Any("args", args))
	err := s.Send(ctx, biz, args, numbers...)
	if err != nil {
		zap.L().Debug("发送短信出现异常", zap.Error(err))
	}

	return err
}
