package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	sredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/service/sms/memory"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	db := initDb()
	rdb := initCache()
	server := initWebServer(rdb)
	c := initUser(db, rdb)
	c.RegisterRoutes(server.Group("/users"))

	server.Run("0.0.0.0:8081")
}

func initWebServer(rdb redis.Cmdable) *gin.Engine {
	server := gin.Default()
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	server.Use(cors.New(cors.Config{
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
	}))
	store, err := sredis.NewStore(16, "tcp", config.Config.Redis.Addr, "", []byte("NwuM65iCW22CiwzIx8t7cmzhAYmBnWUL"), []byte("bOsXTNQzQ1kCAQ9aTWiTtUyuyWEfv5Sf"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store))
	//server.Use(middleware.NewLoginMiddlewareBuild().IgnorePaths("/users/signup", "/users/login").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuild().IgnorePaths("/users/signup", "/users/login", "/users/login_sms/code/send", "/users/login_sms/code/verify").Build())
	return server
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	//测试用
	smsService := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsService)
	c := web.NewUserHandler(svc, codeSvc)
	return c
}

func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
func initCache() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: "",
	})
}
