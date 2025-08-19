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
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/integration/startup"
	"webookProgram/webook/internal/repository/dao/article"
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
	a.db.Exec("TRUNCATE TABLE publish_articles")
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
				var art article.Article
				err := a.db.Where("id=?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
				art.CTime = 0
				art.UTime = 0
				assert.Equal(t, article.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Status:   domain.ArticleUnPublished.ToUint8(),
				}, art)
			},
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
		{
			name: "修改文章-保存成功",
			art: Article{
				Id:      3,
				Title:   "我的标题",
				Content: "我的内容",
			},
			before: func(t *testing.T) {
				art := &article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					CTime:    123,
					UTime:    123,
				}
				err := a.db.Create(art).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := a.db.Where("id=?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
				art.CTime = 0
				art.UTime = 0
				assert.Equal(t, article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Status:   domain.ArticleUnPublished.ToUint8(),
				}, art)
			},
			wantRes: Result[int64]{
				Data: 3,
				Msg:  "OK",
			},
		},
		{
			name: "修改他人文章-保存失败",
			art: Article{
				Id:      4,
				Title:   "我的标题",
				Content: "我的内容",
			},
			before: func(t *testing.T) {
				art := &article.Article{
					Id:       4,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 1234,
					CTime:    123,
					UTime:    123,
					Status:   domain.ArticlePublished.ToUint8(),
				}
				err := a.db.Create(art).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := a.db.Where("id=?", 4).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, article.Article{
					Id:       4,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 1234,
					CTime:    123,
					UTime:    123,
					Status:   domain.ArticlePublished.ToUint8(),
				}, art)
			},
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
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
func (s *ArticleSuite) TestPublish() {
	t := s.T()

	testCases := []struct {
		name string
		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)
		req   Article

		// 预期响应
		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "新建帖子并发表",
			before: func(t *testing.T) {
				// 什么也不需要做
			},
			after: func(t *testing.T) {
				// 验证一下数据
				var art article.Article
				err := s.db.Where("author_id = ?", 123).First(&art).Error
				assert.NoError(t, err)
				// 确保已经生成了主键
				assert.True(t, art.Id > 0)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
				art.CTime = 0
				art.UTime = 0
				art.Id = 0
				assert.Equal(t, article.Article{
					Title:    "hello，你好",
					Content:  "随便试试",
					AuthorId: 123,
					Status:   uint8(domain.ArticlePublished),
				}, art)
				var publishedArt article.PublishArticle
				err = s.db.Where("author_id = ?", 123).First(&publishedArt).Error
				assert.NoError(t, err)
				assert.True(t, publishedArt.Id > 0)
				assert.True(t, publishedArt.CTime > 0)
				assert.True(t, publishedArt.UTime > 0)
				publishedArt.CTime = 0
				publishedArt.UTime = 0
				publishedArt.Id = 0
				assert.Equal(t, article.PublishArticle{
					Article: article.Article{
						Title:    "hello，你好",
						Content:  "随便试试",
						AuthorId: 123,
						Status:   uint8(domain.ArticlePublished),
					},
				}, publishedArt)
			},
			req: Article{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Msg:  "OK",
				Data: 1,
			},
		},
		{
			// 制作库有，但是线上库没有
			name: "更新帖子并新发表",
			before: func(t *testing.T) {
				// 模拟已经存在的帖子
				err := s.db.Create(&article.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					CTime:    456,
					UTime:    234,
					AuthorId: 123,
					Status:   uint8(domain.ArticleUnPublished),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证一下数据
				var art article.Article
				s.db.Where("id = ?", 2).First(&art)
				// 更新时间变了
				assert.True(t, art.UTime > 234)
				art.UTime = 0
				// 创建时间没变
				assert.Equal(t, article.Article{
					Id:       2,
					CTime:    456,
					Status:   uint8(domain.ArticlePublished),
					Content:  "新的内容",
					Title:    "新的标题",
					AuthorId: 123,
				}, art)

				var publishedArt article.PublishArticle
				s.db.Where("id = ?", 2).First(&publishedArt)
				assert.True(t, publishedArt.CTime > 0)
				assert.True(t, publishedArt.UTime > 0)
				publishedArt.CTime = 0
				publishedArt.UTime = 0
				assert.Equal(t, article.PublishArticle{
					Article: article.Article{
						Id:       2,
						Status:   uint8(domain.ArticlePublished),
						Content:  "新的内容",
						Title:    "新的标题",
						AuthorId: 123,
					},
				}, publishedArt)
			},
			req: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Msg:  "OK",
				Data: 2,
			},
		},
		{
			name: "更新帖子，并且重新发表",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					CTime:    456,
					UTime:    234,
					AuthorId: 123,
					Status:   uint8(domain.ArticlePublished),
				}
				err := s.db.Create(&art).Error
				assert.NoError(t, err)
				part := article.PublishArticle{
					Article: art,
				}
				err = s.db.Create(&part).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id = ?", 3).First(&art).Error
				assert.NoError(t, err)
				// 更新时间变了
				assert.True(t, art.UTime > 234)
				art.UTime = 0
				// 创建时间没变
				assert.Equal(t, article.Article{
					Id:       3,
					CTime:    456,
					Status:   uint8(domain.ArticlePublished),
					Content:  "新的内容",
					Title:    "新的标题",
					AuthorId: 123,
				}, art)

				var publishedArt article.PublishArticle
				err = s.db.Where("id = ?", 3).First(&publishedArt).Error
				assert.NoError(t, err)
				assert.True(t, publishedArt.CTime > 0)
				assert.True(t, publishedArt.UTime > 0)
				publishedArt.CTime = 0
				publishedArt.UTime = 0
				assert.Equal(t, article.PublishArticle{
					Article: article.Article{
						Id:       3,
						Status:   uint8(domain.ArticlePublished),
						Content:  "新的内容",
						Title:    "新的标题",
						AuthorId: 123,
					},
				}, publishedArt)
			},
			req: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Msg:  "OK",
				Data: 3,
			},
		},
		{
			name: "更新别人的帖子，并且发表失败",
			before: func(t *testing.T) {
				art := article.Article{
					Id:      4,
					Title:   "我的标题",
					Content: "我的内容",
					CTime:   456,
					UTime:   234,
					// 注意。这个 AuthorID 我们设置为另外一个人的ID
					AuthorId: 789,
				}
				s.db.Create(&art)
				part := article.PublishArticle{
					Article: article.Article{
						Id:       4,
						Title:    "我的标题",
						Content:  "我的内容",
						CTime:    456,
						UTime:    234,
						AuthorId: 789,
					},
				}
				s.db.Create(&part)
			},
			after: func(t *testing.T) {
				// 更新应该是失败了，数据没有发生变化
				var art article.Article
				s.db.Where("id = ?", 4).First(&art)
				assert.Equal(t, "我的标题", art.Title)
				assert.Equal(t, "我的内容", art.Content)
				assert.Equal(t, int64(456), art.CTime)
				assert.Equal(t, int64(234), art.UTime)
				assert.Equal(t, int64(789), art.AuthorId)

				var part article.PublishArticle
				// 数据没有变化
				s.db.Where("id = ?", 4).First(&part)
				assert.Equal(t, "我的标题", part.Title)
				assert.Equal(t, "我的内容", part.Content)
				assert.Equal(t, int64(789), part.AuthorId)
				// 创建时间没变
				assert.Equal(t, int64(456), part.CTime)
				// 更新时间变了
				assert.Equal(t, int64(234), part.UTime)
			},
			req: Article{
				Id:      4,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			// 利用泛型来限定结果必须是 int64
			var result Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
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
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
