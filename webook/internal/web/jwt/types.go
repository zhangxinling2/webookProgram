package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJwtToken(ctx *gin.Context, uid int64, ssid string) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	ExtractToken(ctx *gin.Context) string
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
}
type UserClaims struct {
	jwt.RegisteredClaims
	Id        int64
	UserAgent string
	Ssid      string
}
type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}
