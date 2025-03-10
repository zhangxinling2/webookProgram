package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webookProgram/webook/internal/service"
	svcmock "webookProgram/webook/internal/service/mocks"
)

func TestUserHandler_SignUp(t *testing.T) {

	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#123456",
	"confirmPassword":"hello#123456"
}
`,
			wantCode: 200,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmock.NewMockUserService(ctrl)

				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#123456"
}
`,
			wantCode: http.StatusBadRequest,
			wantBody: "",
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmock.NewMockUserService(ctrl)

				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"123@22qq",
	"password":"hello#123456",
"confirmPassword":"hello#123456"
}
`,
			wantCode: 200,
			wantBody: "邮箱格式不对",
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmock.NewMockUserService(ctrl)

				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#123456",
"confirmPassword":"hello#12345"
}
`,
			wantCode: 200,
			wantBody: "两次密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmock.NewMockUserService(ctrl)

				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello123",
"confirmPassword":"hello123"
}
`,
			wantCode: 200,
			wantBody: "密码必须大于8位，包含数字，特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicateEmail)
				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#123456",
	"confirmPassword":"hello#123456"
}
`,
			wantCode: 200,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("系统异常error"))
				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#123456",
	"confirmPassword":"hello#123456"
}
`,
			wantCode: 200,
			wantBody: "系统异常",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()

			h := NewUserHandler(tt.mock(ctrl), nil)
			h.RegisterRoutes(server.Group("/users"))
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tt.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, tt.wantBody, resp.Body.String())
		})
	}
}
