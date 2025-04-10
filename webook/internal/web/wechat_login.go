package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
	"webookProgram/webook/internal/service"
	"webookProgram/webook/internal/service/oauth2"
)

type OAuth2WechatHandler struct {
	svc       service.UserService
	OAuth2Svc oauth2.Service
	jwtHandler
	stateKey []byte
	config   Config
}

func NewOAuth2WechatHandler(svc service.UserService, OAuth2Svc oauth2.Service, config Config) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:       svc,
		OAuth2Svc: OAuth2Svc,
		stateKey:  []byte("6uZhFEhonyX0JalbKDkarQMRpzLwuS3A"),
		config:    config,
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	})
	tokenString, err := token.SignedString(w.stateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.SetCookie("jwt-state", tokenString, 600, "/oauth2/wechat/callback", "", w.config.Secure, w.config.HttpOnly)
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Data: urlPattern,
	})
}
func (w *OAuth2WechatHandler) VerifyCode(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "登录失败",
		})
		return
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return w.stateKey, nil
	})
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "登录失败",
		})
		return
	}
	if sc.State != state {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "登录失败",
		})
		return
	}
	info, err := w.OAuth2Svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	u, err := w.svc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	err = w.setJwtToken(ctx, u.Id)
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
type Config struct {
	Secure   bool
	HttpOnly bool
}
