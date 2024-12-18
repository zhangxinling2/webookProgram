package web

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func (h *UserHandler) signup(ctx *gin.Context) {

}
func (h *UserHandler) profile(ctx *gin.Context) {

}
func (h *UserHandler) edit(ctx *gin.Context) {

}
func (h *UserHandler) login(ctx *gin.Context) {

}
func (h *UserHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/signup", h.signup)
	group.POST("/edit", h.edit)
	group.POST("/login", h.login)
	group.GET("/profile", h.profile)
}
