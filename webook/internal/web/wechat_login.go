package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
	"webookProgram/webook/internal/service"
	"webookProgram/webook/internal/service/oauth2"
	jwt2 "webookProgram/webook/internal/web/jwt"
)

type OAuth2WechatHandler struct {
	svc       service.UserService
	OAuth2Svc oauth2.Service
	jwt2.Handler
	stateKey []byte
	config   WechatHandlerConfig
}

func NewOAuth2WechatHandler(svc service.UserService, OAuth2Svc oauth2.Service, config WechatHandlerConfig, handler jwt2.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:       svc,
		OAuth2Svc: OAuth2Svc,
		stateKey:  []byte("6uZhFEhonyX0JalbKDkarQMRpzLwuS3S"),
		config:    config,
		Handler:   handler,
	}
}
func (w *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := shortuuid.New()
	urlPattern, err := w.OAuth2Svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if err = w.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Data: urlPattern,
	})
}

func (w *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	})
	tokenString, err := token.SignedString(w.stateKey)
	if err != nil {
		return err
	}
	ctx.SetCookie("jwt-state", tokenString, 600, "/oauth2/wechat/callback", "", w.config.Secure, w.config.HttpOnly)
	return nil
}
func (w *OAuth2WechatHandler) VerifyCode(ctx *gin.Context) {
	code := ctx.Query("code")
	err := w.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	info, err := w.OAuth2Svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	u, err := w.svc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = w.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (w *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("%w, 无法获得 Cookie", err)
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return w.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("%w, token已经过期", err)
	}
	if sc.State != state {
		return errors.New("state不相等")
	}
	return err
}
func NewWechatHandler() *OAuth2WechatHandler {
	return &OAuth2WechatHandler{}
}
func (w *OAuth2WechatHandler) RegisterGroup(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", w.AuthURL)
	g.Any("/verifycode", w.VerifyCode)
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
type WechatHandlerConfig struct {
	Secure   bool
	HttpOnly bool
}
