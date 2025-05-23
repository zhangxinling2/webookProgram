package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webookProgram/webook/internal/web"
	jwt2 "webookProgram/webook/internal/web/jwt"
	"webookProgram/webook/internal/web/middleware"
	ratelimit2 "webookProgram/webook/pkg/ginx/middlewares/ratelimit"
	"webookProgram/webook/pkg/ratelimit"
)

func InitEngine(mdls []gin.HandlerFunc, hdl *web.UserHandler, whdl *web.OAuth2WechatHandler, ahdl *web.ArticleHandler) *gin.Engine {
	//store, err := sredis.NewStore(16, "tcp", config.WechatHandlerConfig.Redis.Addr, "", []byte("NwuM65iCW22CiwzIx8t7cmzhAYmBnWUL"), []byte("bOsXTNQzQ1kCAQ9aTWiTtUyuyWEfv5Sf"))
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("ssid", store))
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server.Group("/users"))
	ahdl.RegisterRoutes(server.Group("/article"))
	whdl.RegisterGroup(server)
	return server
}
func InitSlideWindowLimit(redisClient redis.Cmdable) ratelimit.Limiter {
	return ratelimit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 100)
}
func InitMiddlewares(rateLimit ratelimit.Limiter, jwtHdl jwt2.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuild(jwtHdl).IgnorePaths("/users/signup", "/users/login", "/users/login_sms/code/send", "/users/login_sms", "/users/refresh_token").Build(),
		ratelimit2.NewBuilder(rateLimit).Build(),
	}

}
func NewRedisJwtHandler(redisClient redis.Cmdable) jwt2.Handler {
	return jwt2.NewRedisJwtHandler(redisClient)
}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{"PUT", "Post","GET"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://localhost") || strings.Contains(origin, "webook.com") {
				//开发环境
				return true
			}
			return strings.Contains(origin, "公司域名")
		},
		MaxAge: 12 * time.Hour,
	})
}
