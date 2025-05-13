package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/service"
	"webookProgram/webook/internal/web/jwt"
	"webookProgram/webook/pkg/logger"
)

var _ handler = &ArticleHandler{}

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}
func (a *ArticleHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	type Article struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var art Article
	if err := ctx.Bind(&art); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c := ctx.MustGet("claims")
	uc, ok := c.(*jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("未找到用户session信息")
		return
	}
	id, err := a.svc.Save(ctx, domain.Article{
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "OK",
		Data: id,
	})
}
