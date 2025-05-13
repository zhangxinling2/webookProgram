package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"webookProgram/webook/internal/web"
)

type ArticleSuite struct {
	suite.Suite
	server *gin.Engine
}

func (a *ArticleSuite) SetupSuite() {

}
func (a *ArticleSuite) TearDownSuite() {

}
func (a *ArticleSuite) TearDownTest() {

}
func (a *ArticleSuite) TestEdit(t *testing.T) {
	testCases := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantRes Result[int64]
		art     Article
	}{
		{},
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
			a.server.ServeHTTP(resp, req)
			var res web.Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)
		})
	}
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
