package freecache

import (
	"context"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"sync"
	"webookProgram/webook/internal/repository/cache"
)

var (
	ErrSetCodeSendTooMany = errors.New("发送太频繁")
	ErrCodeVerifyTooMany  = errors.New("验证次数太多")
	ErrUnknown            = errors.New("未知错误")
)

type CodeFreeCache struct {
	cache *freecache.Cache
	lock  sync.Mutex
}

func NewCodeCache(cache *freecache.Cache, lock sync.Mutex) cache.CodeCache {
	return &CodeFreeCache{cache: cache}
}
func (c *CodeFreeCache) Set(ctx context.Context, biz, phone, code string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	key := []byte(c.key(biz, phone))
	time, err := c.cache.TTL(key)
	if err == freecache.ErrNotFound || time < 540 {
		err = c.cache.Set(key, []byte(code), 600)
		if err != nil {
			return err
		}
		keyCnt := append(key, []byte("Cnt")...)
		err = c.cache.Set(keyCnt, []byte("3"), 600)
		if err != nil {
			return err
		}
	} else if time > 540 {
		return ErrSetCodeSendTooMany
	} else {
		return ErrUnknown
	}
	return nil
}

func (c *CodeFreeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	key := []byte(c.key(biz, phone))
	val, err := c.cache.Get(key)
	if err != nil {
		return false, err
	}
	keyCnt := append(key, []byte("Cnt")...)
	cntTmpVal, err := c.cache.Get(keyCnt)
	if err != nil {
		return false, err
	}
	cntVal := cntTmpVal[0] - '0'
	if cntVal <= 0 {
		return false, ErrCodeVerifyTooMany
	} else if string(val) == inputCode {
		err = c.cache.Set(keyCnt, []byte("-1"), 600)
		return true, err
	} else {
		err = c.cache.Set(keyCnt, []byte(string(cntVal-1)), 600)
		return false, err
	}

}
func (c *CodeFreeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
