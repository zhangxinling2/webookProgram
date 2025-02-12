package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

	userIdKey = "userId"
	bizLogin  = "login"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:         svc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}
func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
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
	c, ok := ctx.Get("claims")
	if !ok {
		//可以考虑监控这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	u, err := h.svc.Profile(ctx, claims.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
	}
	ctx.JSON(http.StatusOK, u)
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

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid: user.Id}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("6uZhFEhonyX0JalbKDkarQMRpzLwuS3N"))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.Header("x-jwt-token", tokenStr)
}
func (h *UserHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/signup", h.SignUp)
	group.POST("/edit", h.edit)
	//group.POST("/login", h.login)
	group.POST("/login", h.loginJWT)
	//group.GET("/profile", h.profile)
	group.GET("/profile", h.profileJWT)
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
