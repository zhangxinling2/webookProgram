package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	db := initDb()
	server := initWebServer()
	c := initUser(db)
	c.RegisterRoutes(server.Group("/users"))
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{"PUT", "Post","GET"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://localhost") {
				//开发环境
				return true
			}
			return strings.Contains(origin, "公司域名")
		},
		MaxAge: 12 * time.Hour,
	}))
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("NwuM65iCW22CiwzIx8t7cmzhAYmBnWUL"), []byte("bOsXTNQzQ1kCAQ9aTWiTtUyuyWEfv5Sf"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store))
	//server.Use(middleware.NewLoginMiddlewareBuild().IgnorePaths("/users/signup", "/users/login").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuild().IgnorePaths("/users/signup", "/users/login").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	c := web.NewUserHandler(svc)
	return c
}

func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
