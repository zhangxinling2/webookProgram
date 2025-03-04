package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webookProgram/webook/internal/web"
	"webookProgram/webook/internal/web/middleware"
	"webookProgram/webook/pkg/ginx/middlewares/ratelimit"
)

func InitEngine(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	//store, err := sredis.NewStore(16, "tcp", config.Config.Redis.Addr, "", []byte("NwuM65iCW22CiwzIx8t7cmzhAYmBnWUL"), []byte("bOsXTNQzQ1kCAQ9aTWiTtUyuyWEfv5Sf"))
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("ssid", store))
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server.Group("/users"))
	return server
}
func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuild().IgnorePaths("/users/signup", "/users/login", "/users/login_sms/code/send", "/users/login_sms").Build(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}

}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{"PUT", "Post","GET"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
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
