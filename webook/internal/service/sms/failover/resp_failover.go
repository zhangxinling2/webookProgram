package failover

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository"
	"webookProgram/webook/internal/service/sms"
	"webookProgram/webook/pkg/logger"
)

type RespTimeGrowSMSService struct {
	svcs        []sms.Service
	baseTime    time.Duration
	curInx      int32
	windowsSize int
	svcWindow   []*timeWindow
	repo        repository.AsyncSmsRepository
	l           logger.LoggerV1
}
type timeWindow struct {
	times   []time.Duration
	sumTime time.Duration
	cursor  int
}

func NewRespTimeGrowSMSService(svcs []sms.Service, curInx int32, windowSize int, repo repository.AsyncSmsRepository, l logger.LoggerV1) sms.Service {
	windows := make([]*timeWindow, len(svcs))
	for i, _ := range windows {
		windows[i] = &timeWindow{
			times: make([]time.Duration, windowSize),
		}
	}
	return &RespTimeGrowSMSService{svcs: svcs, curInx: curInx, windowsSize: windowSize, svcWindow: windows, repo: repo, l: l}
}

// 异步发送消息，最简单的抢占式调度
func (f *RespTimeGrowSMSService) StartAsyncCircle() {
	for {
		f.AsyncSend()
	}
}
func (f *RespTimeGrowSMSService) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	sms, err := f.repo.PrepareWaitingSms(ctx)
	cancel()
	switch err {
	case nil:
		idx := atomic.LoadInt32(&f.curInx)
		err = f.svcs[idx].Send(ctx, sms.Biz, sms.Args, sms.Numbers...)
		if err != nil {
			f.l.Error("执行异步发送失败", logger.Error(err), logger.Int64("id", sms.Id))
		}
		res := err == nil
		err = f.repo.ReportScheduleResult(ctx, sms.Id, res)
		if err != nil {
			f.l.Error("异步发送成功，报告结果失败失败", logger.Error(err), logger.Bool("res", res), logger.Int64("id", sms.Id))
		}
	case repository.ErrNoSms:
		time.Sleep(1000)
	default:
		f.l.Error("获取异步发送数据失败", logger.Error(err))
	}
}

func (f *RespTimeGrowSMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&f.curInx)
	window := f.svcWindow[idx]
	now := time.Now()
	err := f.svcs[idx].Send(ctx, biz, args, numbers...)
	duration := time.Since(now)

	window.times[window.cursor] = duration
	window.cursor = (window.cursor + 1) % f.windowsSize
	window.sumTime = window.sumTime + duration
	curAvg := window.sumTime / time.Duration(len(window.times))
	if f.baseTime == 0 {
		f.baseTime = curAvg
		return err
	}
	if curAvg > f.baseTime*120/100 {
		newIdx := int(f.curInx+1) % len(f.svcs)
		atomic.CompareAndSwapInt32(&f.curInx, idx, int32(newIdx))
		smsDomain := domain.AsyncSms{
			Biz:      biz,
			Args:     args,
			Numbers:  numbers,
			RetryMax: 3,
		}
		err = f.repo.Store(ctx, smsDomain)
		if err != nil {
			return err
		}
	}
	return errors.New("全部服务商都失败")
}
