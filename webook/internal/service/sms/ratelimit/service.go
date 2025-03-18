package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"webookProgram/webook/internal/service/sms"
	"webookProgram/webook/pkg/ratelimit"
)

type RateLimiterSmsService struct {
	svc   sms.Service
	limit ratelimit.Limiter
}

func NewRateLimiterService(svc sms.Service, limit ratelimit.Limiter) sms.Service {
	return &RateLimiterSmsService{svc: svc, limit: limit}
}

const key = "sms_limit"

func (s *RateLimiterSmsService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	ok, err := s.limit.Limit(ctx, key)
	if err != nil {
		return fmt.Errorf("短信限流错误:%s", err)
	}
	if ok {
		return errors.New("短信限流")
	}
	return s.svc.Send(ctx, biz, args, numbers...)
}
