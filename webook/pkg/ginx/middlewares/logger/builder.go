package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type MiddlewareBuild struct {
	// body可能很大，所以要小心
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuild(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuild {
	return &MiddlewareBuild{
		loggerFunc: fn,
	}
}
func (b *MiddlewareBuild) AllowReqBody() *MiddlewareBuild {
	b.allowReqBody = true
	return b
}
func (b *MiddlewareBuild) AllowRespBody() *MiddlewareBuild {
	b.allowRespBody = true
	return b
}
func (b *MiddlewareBuild) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		url := c.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: c.Request.Method,
			Url:    url,
		}
		if b.allowReqBody && c.Request.Body != nil {
			body, _ := c.GetRawData()
			//要放回来，因为body是一个ReadCloser，读完就没了
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			// 其实是一个很消耗CPU 和 内存的操作，因为会引起复制
			al.ReqBody = string(body)
		}
		defer func() {
			al.Duration = time.Since(start)
			if b.allowRespBody && c.Request.Body != nil {
				body, _ := io.ReadAll(c.Request.Body)
				//要放回来，因为body是一个ReadCloser，读完就没了
				c.Request.Body = io.NopCloser(bytes.NewReader(body))
				// 其实是一个很消耗CPU 和 内存的操作，因为会引起复制
				al.ReqBody = string(body)
			}
			b.loggerFunc(c, al)
		}()
		c.Next()

	}
}

type AccessLog struct {
	Method   string
	Url      string
	ReqBody  string
	RespBody string
	Duration time.Duration
}
