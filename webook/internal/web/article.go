package web

import "github.com/gin-gonic/gin"

var _ handler = &ArticleHandler{}

type ArticleHandler struct {
}

func (a *ArticleHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {

}
