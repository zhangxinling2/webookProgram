package main

import (
	"github.com/gin-gonic/gin"
	"webook/internal/web"
)

func main() {
	server := gin.Default()
	c := &web.UserHandler{}
	c.RegisterRoutes(server.Group("/users"))
	captcha := web.NewCaptchaHandler()
	captcha.RegisterRoutes(server.Group("/captcha"))
	server.Run("0.0.0.0:3000")
}
