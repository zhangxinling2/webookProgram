package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

var (
	AtKey = []byte("6uZhFEhonyX0JalbKDkarQMRpzLwuS3A")
	RtKey = []byte("6uZhFEhonyX0JalbKDkarQMRpzLwuS3R")
)

func NewRedisJwtHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{

		cmd: cmd,
	}

}
func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	claims, ok := ctx.Get("claims")
	if !ok {
		return errors.New("claims not found")
	}
	userClaims := claims.(*UserClaims)
	err := h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", userClaims.Ssid), "", time.Hour*24*7).Err()
	if err != nil {
		return err
	}
	return nil
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	return err
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJwtToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return nil
}
func (h *RedisJWTHandler) SetJwtToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
func (h *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-Refresh-token", tokenStr)
	return nil
}
func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	if tokenHeader == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return segs[1]
}
