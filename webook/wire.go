//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webookProgram/webook/internal/repository"
	"webookProgram/webook/internal/repository/cache"
	redis2 "webookProgram/webook/internal/repository/cache/redis"
	"webookProgram/webook/internal/repository/dao"
	"webookProgram/webook/internal/service"
	"webookProgram/webook/internal/service/sms"
	"webookProgram/webook/internal/service/sms/memory"
	"webookProgram/webook/internal/web"
	"webookProgram/webook/ioc"
)

var CodeCacheSet = wire.NewSet(redis2.NewCodeCache,
	wire.Bind(new(cache.CodeCache), new(*redis2.CodeRedisCache)))
var SmsServiceSet = wire.NewSet(memory.NewService,
	wire.Bind(new(sms.Service), new(*memory.Service)))

func InitWebServer() *gin.Engine {
	wire.Build(
		//基础第三方依赖
		ioc.InitDb, ioc.InitCache,
		ioc.InitLogger,
		//DAO
		dao.NewUserDAO,
		cache.NewUserCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		ioc.NewWechatHandlerConfig,
		ioc.InitSMSService,
		ioc.InitOAuth2WechatService,
		ioc.NewRedisJwtHandler,
		redis2.NewCodeCache,
		web.NewOAuth2WechatHandler,
		web.NewUserHandler,
		ioc.InitSlideWindowLimit,
		ioc.InitMiddlewares,
		ioc.InitEngine,
	)
	return new(gin.Engine)
}
