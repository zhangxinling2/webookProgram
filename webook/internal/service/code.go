package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

const codeTplId = "1877556"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{repo: repo,
		smsSvc: smsSvc}
}
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	//如果这里出错可以加上重试，但是重试要在Service里重试，也就是传一个可以重试的服务，不能删掉这个验证码因为可能是超时的错误，无法判断是否发出
	return err
}
func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	//不够6位的加上前导0
	return fmt.Sprintf("%6d", num)
}
func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}
