package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/service/article"
	"webookProgram/webook/internal/web/jwt"
	"webookProgram/webook/pkg/logger"
)

var _ handler = &ArticleHandler{}

type ArticleHandler struct {
	svc article.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc article.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}
func (a *ArticleHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/edit", a.Edit)
	group.POST("/publish", a.Publish)
	group.POST("/withdraw", a.Withdraw)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
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
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: uc.Id,
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

func (a *ArticleHandler) Publish(ctx *gin.Context) {
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
	id, err := a.svc.Publish(ctx, domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: uc.Id,
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

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
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
	err := a.svc.Withdraw(ctx, req.Id, uc.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("撤回帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "OK",
	})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
