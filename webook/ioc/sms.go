package ioc

import (
	"webookProgram/webook/internal/service/sms"
	"webookProgram/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
