package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
)

type Service struct {
	appId     *string
	signature *string
	client    *sms.Client
}

func NewService(client *sms.Client, appId string, signature string) *Service {
	return &Service{
		appId:     ekit.ToPtr[string](appId),
		signature: ekit.ToPtr[string](signature),
		client:    client,
	}
}
func (s *Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signature
	req.TemplateId = ekit.ToPtr[string](biz)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSms(req)
	zap.L().Debug("发送短信", zap.Any("req", req), zap.Any("resp", resp), zap.Error(err))
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) == "Ok" {
			return fmt.Errorf("发送短信失败 %s,%s", status.Code, status.Message)
		}
	}
	return nil
}
func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
