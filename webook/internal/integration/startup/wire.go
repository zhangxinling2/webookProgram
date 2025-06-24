//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webookProgram/webook/internal/repository"
	article2 "webookProgram/webook/internal/repository/article"
	"webookProgram/webook/internal/repository/cache"
	redis2 "webookProgram/webook/internal/repository/cache/redis"
	"webookProgram/webook/internal/repository/dao"
	article3 "webookProgram/webook/internal/repository/dao/article"
	"webookProgram/webook/internal/service"
	"webookProgram/webook/internal/service/article"
	"webookProgram/webook/internal/service/sms"
	"webookProgram/webook/internal/service/sms/memory"
	"webookProgram/webook/internal/web"
	"webookProgram/webook/ioc"
)

var thirdProvider = wire.NewSet(InitTestDb, InitCache, InitLogger)
var CodeCacheSet = wire.NewSet(redis2.NewCodeCache,
	wire.Bind(new(cache.CodeCache), new(*redis2.CodeRedisCache)))
var SmsServiceSet = wire.NewSet(memory.NewService,
	wire.Bind(new(sms.Service), new(*memory.Service)))

func InitWebServer() *gin.Engine {
	wire.Build(
		//基础第三方依赖
		thirdProvider,
		//DAO
		dao.NewUserDAO,
		article3.NewArticleDAO,
		cache.NewUserCache,
		redis2.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		article2.NewArticleRepository,
		service.NewUserService,
		service.NewCodeService,
		article.NewArticleService,
		ioc.NewWechatHandlerConfig,
		ioc.InitSMSService,
		ioc.InitOAuth2WechatService,
		ioc.NewRedisJwtHandler,
		web.NewOAuth2WechatHandler,
		web.NewUserHandler,
		web.NewArticleHandler,
		ioc.InitSlideWindowLimit,
		ioc.InitMiddlewares,
		ioc.InitEngine,
	)
	return new(gin.Engine)
}
func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdProvider,
		article3.NewArticleDAO,
		article2.NewArticleRepository,
		article.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}
