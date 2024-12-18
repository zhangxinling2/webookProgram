package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/utils"
)

type CaptchaHandler struct {
	captcha *utils.StringCaptcha
}
type CaptchaResult struct {
	id  string
	ans string
}

func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{
		captcha: utils.NewCaptcha(),
	}
}
func (c *CaptchaHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/getCaptcha", c.Get)
	group.POST("/verify", c.Verify)
}
func (c *CaptchaHandler) Get(ctx *gin.Context) {
	id, base64, ans := c.captcha.Generate()
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"base64":    base64,
		"answer":    ans,
	})
}
func (c *CaptchaHandler) Verify(ctx *gin.Context) {
	var cr CaptchaResult
	ctx.Bind(&cr)
	var res = c.captcha.Verify(cr.id, cr.ans)
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": cr.id,
		"answer":    cr.ans,
		"result":    res,
	})
}
