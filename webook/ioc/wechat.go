package ioc

import (
	"webookProgram/webook/internal/service/oauth2"
	"webookProgram/webook/internal/web"
)

func InitOAuth2WechatService() oauth2.Service {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("WECHAT_APP_ID env variable is not set")
	//}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_Secret")
	//if !ok {
	//	panic("WECHAT_APP_Secret env variable is not set")
	//}
	return oauth2.NewService("appid", "appSecret")
}
func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
