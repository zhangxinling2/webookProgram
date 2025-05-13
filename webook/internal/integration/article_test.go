package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"webookProgram/webook/internal/integration/startup"
	"webookProgram/webook/internal/repository/dao"
	"webookProgram/webook/internal/web/jwt"
)

type ArticleSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (a *ArticleSuite) SetupSuite() {
	a.server = gin.Default()
	a.server.Use(func(context *gin.Context) {
		context.Set("claims", &jwt.UserClaims{
			Uid: 123,
		})
	})
	a.db = startup.InitTestDb()
	ahdl := startup.InitArticleHandler()
	ahdl.RegisterRoutes(a.server.Group("/articles"))
}
func (a *ArticleSuite) TearDownSuite() {

}
func (a *ArticleSuite) TearDownTest() {
	a.db.Exec("TRUNCATE TABLE articles")
}
func (a *ArticleSuite) TestEdit() {
	t := a.T()
	testCases := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantRes Result[int64]
		art     Article
	}{
		{
			name: "创建文章-保存成功",
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				err := a.db.Where("id=?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
				art.CTime = 0
				art.UTime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
				}, art)
			},
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			a.server.ServeHTTP(resp, req)
			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)
		})
	}
}
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleSuite{})
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
