package startup

import (
	"webookProgram/webook/internal/service/sms"
	"webookProgram/webook/internal/service/sms/memory"
	smslimiter "webookProgram/webook/internal/service/sms/ratelimit"
	"webookProgram/webook/pkg/ratelimit"
)

func InitSMSService(limiter ratelimit.Limiter) sms.Service {
	svc := memory.NewService()
	return smslimiter.NewRateLimiterService(svc, limiter)
}
