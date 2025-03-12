package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webookProgram/webook/internal/web"
	"webookProgram/webook/ioc"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitCache()
	testCases := []struct {
		name    string
		reqBody string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantRes web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:12311111111").Result()
				cancel()
				require.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
		"phone":"12311111111"
}
`,
			wantRes: web.Result{
				Code: 0,
				Msg:  "发送成功",
				Data: nil,
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.Set(ctx, "phone_code:login:12311111111", "123456", 0).Result()
				cancel()
				require.NoError(t, err)
				assert.Equal(t, val, "OK")
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:12311111111").Result()
				cancel()
				require.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
		"phone":"12311111111"
}
`,
			wantRes: web.Result{
				Code: 5,
				Msg:  "系统错误",
				Data: nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			var res web.Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)
		})
	}
}
