package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/service"
	jwt2 "webookProgram/webook/internal/web/jwt"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

	userIdKey = "userId"
	bizLogin  = "login"
)

var _ handler = &UserHandler{}

type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	jwt2.Handler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, handler jwt2.Handler) *UserHandler {
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		Handler:     handler,
	}
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	ok, err = h.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字，特殊字符")
		return
	}
	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	// 数据库操作
}
func (h *UserHandler) profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	id := sess.Get("userId")
	u, err := h.svc.Profile(ctx, id.(int64))
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
	}
	ctx.JSON(http.StatusOK, u)
}
func (h *UserHandler) profileJWT(ctx *gin.Context) {
	type Profile struct {
		Email    string
		Phone    string
		Nickname string
		Birth    string
		About    string
	}
	c, ok := ctx.Get("claims")
	if !ok {
		//可以考虑监控这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	claims, ok := c.(*jwt2.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	u, err := h.svc.Profile(ctx, claims.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
	}
	ctx.JSON(http.StatusOK, Profile{
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.NickName,
		Birth:    u.Birth.Format(time.DateOnly),
		About:    u.Introduction,
	})
}
func (h *UserHandler) edit(ctx *gin.Context) {
	type EditReq struct {
		NickName     string `json:"nickName"`
		Birth        string `json:"birth"`
		Introduction string `json:"introduction"`
	}
	var req EditReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	if req.NickName == "" {
		ctx.JSON(http.StatusOK, "昵称不能为空")
		return
	}
	if len(req.Introduction) > 1024 {
		ctx.JSON(http.StatusOK, "关于我过长")
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birth)
	if err != nil {
		// 也就是说，我们其实并没有直接校验具体的格式
		// 而是如果你能转化过来，那就说明没问题
		ctx.JSON(http.StatusOK, "日期格式不对")
		return
	}
	sess := sessions.Default(ctx)
	id := sess.Get("userId")
	err = h.svc.Edit(ctx, domain.User{Id: id.(int64),
		NickName:     req.NickName,
		Birth:        birthday,
		Introduction: req.Introduction,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
	}
	ctx.JSON(http.StatusOK, "成功")
}
func (h *UserHandler) login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.String(http.StatusOK, "登陆成功")
	return
}
func (h *UserHandler) VerifyCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("校验验证码出错", zap.Error(err))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证码有误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码正确",
	})
	//userId从哪来？
	u, err := h.svc.FindOrCreate(ctx, req.Phone)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	if err = h.SetLoginToken(ctx, u.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
}
func (h *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "发送成功",
			Data: nil,
		})
	case service.ErrSetCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后重试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
	}
}
func (h *UserHandler) loginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	err = h.SetLoginToken(ctx, user.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
}
func (h *UserHandler) refreshToken(ctx *gin.Context) {
	tokenStr := h.ExtractToken(ctx)
	if tokenStr == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var rc jwt2.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return jwt2.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = h.SetJwtToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		zap.L().Error("设置 JWT token 出现异常", zap.Error(err))
	}
}
func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登录失败",
			Data: nil,
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录成功",
	})
}
func (h *UserHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/signup", h.SignUp)
	group.POST("/edit", h.edit)
	//group.POST("/login", h.login)
	group.POST("/login", h.loginJWT)
	//group.GET("/profile", h.profile)
	group.GET("/profile", h.profileJWT)
	group.POST("/login_sms/code/send", h.SendLoginSMSCode)
	group.POST("/login_sms", h.VerifyCode)
	group.POST("/refresh_token", h.refreshToken)
	group.POST("/logout", h.LogoutJWT)
}
