package failover

import (
	"context"
	"sync/atomic"
	"webookProgram/webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	svcs      []sms.Service
	idx       int32
	cnt       int32
	threshold int32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, idx int32, cnt int32, threshold int32) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{svcs: svcs, idx: idx, cnt: cnt, threshold: threshold}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		//如果cas操作失败说明别人切换了
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
	}
	idx = atomic.LoadInt32(&t.idx)
	err := t.svcs[idx].Send(ctx, tplId, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		return err
	}
}
