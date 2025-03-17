package failover

import (
	"context"
	"errors"
	"sync/atomic"
	"webookProgram/webook/internal/service/sms"
)

type FailOverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailOverSMSService(svcs []sms.Service) sms.Service {
	return &FailOverSMSService{svcs: svcs}
}

func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		//输出日志，做好监控，意味着全部失败
		return err
	}
	return errors.New("全部服务商都失败")
}
func (f *FailOverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			//服务超时了
			return err
		}
	}
	return errors.New("全部服务商失败")
}
