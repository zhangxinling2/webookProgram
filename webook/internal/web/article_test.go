package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/service/article"
	artmock "webookProgram/webook/internal/service/mocks/article"
	"webookProgram/webook/internal/web/jwt"
	"webookProgram/webook/pkg/logger"
)

func TestArticleHandler_Publish(t *testing.T) {

	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) article.ArticleService
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "新建帖子，发表成功",
			mock: func(ctrl *gomock.Controller) article.ArticleService {
				artSvc := artmock.NewMockArticleService(ctrl)
				artSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return artSvc
			},
			reqBody: `
{
	"title":"我的标题",
	"content":"我的内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				Msg:  "OK",
				Data: float64(1),
			},
		},
		{
			name: "发表失败",
			mock: func(ctrl *gomock.Controller) article.ArticleService {
				artSvc := artmock.NewMockArticleService(ctrl)
				artSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("发表失败"))
				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return artSvc
			},
			reqBody: `
{
	"title":"我的标题",
	"content":"我的内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "Bind 失败",
			mock: func(ctrl *gomock.Controller) article.ArticleService {
				artSvc := artmock.NewMockArticleService(ctrl)
				return artSvc
			},
			reqBody: `
{
	"title":"我的标题"
}
`,
			wantCode: 400,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "修改帖子，发表成功",
			mock: func(ctrl *gomock.Controller) article.ArticleService {
				artSvc := artmock.NewMockArticleService(ctrl)
				artSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "修改标题",
					Content: "修改内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				//codeSvc := svcmock.NewMockCodeService(ctrl)
				return artSvc
			},
			reqBody: `
{
	"id":1,
	"title":"修改标题",
	"content":"修改内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				Msg:  "OK",
				Data: float64(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("claims", &jwt.UserClaims{
					Id: 123,
				})
			})
			h := NewArticleHandler(tt.mock(ctrl), logger.NewNoOpLogger())
			h.RegisterRoutes(server.Group("/article"))
			req, err := http.NewRequest(http.MethodPost, "/article/publish", bytes.NewBuffer([]byte(tt.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantCode, resp.Code)
			var res Result
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantRes, res)
		})
	}
}
