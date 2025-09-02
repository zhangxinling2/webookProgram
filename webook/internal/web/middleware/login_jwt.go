package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	jwt2 "webookProgram/webook/internal/web/jwt"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
	jwt2.Handler
}

func NewLoginJWTMiddlewareBuild(handler jwt2.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: handler,
	}
}
func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path ...string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path...)
	return l
}
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		tokenStr := l.ExtractToken(ctx)
		claims := &jwt2.UserClaims{}
		//一定要传claim指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("6uZhFEhonyX0JalbKDkarQMRpzLwuS3N"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid || claims.Id == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			if claims.UserAgent != ctx.Request.UserAgent() {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Minute))
		//	tokenStr, err = token.SignedString([]byte("6uZhFEhonyX0JalbKDkarQMRpzLwuS3N"))
		//	if err != nil {
		//
		//	}
		//	ctx.Header("x-jwt-token", tokenStr)
		//}
		//以复用结果
		ctx.Set("claims", claims)
	}
}
