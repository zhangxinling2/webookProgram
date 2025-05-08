package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webookProgram/webook/internal/web/jwt"
)

func WrapWithoutClaims[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}
func WrapReq[T any](fn func(ctx *gin.Context, req T, uc jwt.UserClaims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			return
		}
		c, ok := ctx.Get("claims")
		if !ok {
			//可以考虑监控这里
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := c.(jwt.UserClaims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(ctx, req, claims)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
