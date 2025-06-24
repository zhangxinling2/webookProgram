package article

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/article"
	articlemock "webookProgram/webook/internal/repository/article/mocks/article"
	"webookProgram/webook/pkg/logger"
)

func Test_articleServiceV1_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository)
		art     domain.Article
		wantId  int64
		wantErr error
	}{
		{
			name: "新建发表成功",
			art: domain.Article{
				Title:   "新建发表",
				Content: "新建发表",
				Author: domain.Author{
					Id:   123,
					Name: "123",
				},
			},
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := articlemock.NewMockArticleAuthorRepository(ctrl)
				reader := articlemock.NewMockArticleReaderRepository(ctrl)
				ctx := context.Background()
				author.EXPECT().Create(ctx, domain.Article{
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id:   123,
						Name: "123",
					},
				}).Return(int64(1), nil)
				reader.EXPECT().Save(ctx, domain.Article{
					Id:      1,
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id:   123,
						Name: "123",
					},
				}).Return(nil)
				return author, reader
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "修改发表成功",
			art: domain.Article{
				Id:      1,
				Title:   "修改发表",
				Content: "修改发表",
				Author: domain.Author{
					Id:   123,
					Name: "123",
				},
			},
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := articlemock.NewMockArticleAuthorRepository(ctrl)
				reader := articlemock.NewMockArticleReaderRepository(ctrl)
				ctx := context.Background()
				author.EXPECT().Update(ctx, domain.Article{
					Id:      1,
					Title:   "修改发表",
					Content: "修改发表",
					Author: domain.Author{
						Id:   123,
						Name: "123",
					},
				}).Return(int64(1), nil)
				reader.EXPECT().Save(ctx, domain.Article{
					Id:      1,
					Title:   "修改发表",
					Content: "修改发表",
					Author: domain.Author{
						Id:   123,
						Name: "123",
					},
				}).Return(nil)
				return author, reader
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "制作库成功，线上库失败",
			art: domain.Article{
				Id:      1,
				Title:   "修改发表",
				Content: "修改发表",
				Author: domain.Author{
					Id:   123,
					Name: "123",
				},
			},
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := articlemock.NewMockArticleAuthorRepository(ctrl)
				reader := articlemock.NewMockArticleReaderRepository(ctrl)
				ctx := context.Background()
				author.EXPECT().Update(ctx, domain.Article{
					Id:      1,
					Title:   "修改发表",
					Content: "修改发表",
					Author: domain.Author{
						Id:   123,
						Name: "123",
					},
				}).Return(int64(1), nil)
				reader.EXPECT().Save(ctx, domain.Article{
					Id:      1,
					Title:   "修改发表",
					Content: "修改发表",
					Author: domain.Author{
						Id:   123,
						Name: "123",
					},
				}).Return(errors.New("线上库失败")).MaxTimes(3)
				return author, reader
			},
			wantId:  1,
			wantErr: errors.New("线上库失败"),
		},
		{
			name: "制作库新建失败",
			art: domain.Article{
				Title:   "新建发表",
				Content: "新建发表",
				Author: domain.Author{
					Id:   123,
					Name: "123",
				},
			},
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := articlemock.NewMockArticleAuthorRepository(ctrl)
				reader := articlemock.NewMockArticleReaderRepository(ctrl)
				ctx := context.Background()
				author.EXPECT().Create(ctx, domain.Article{
					Title:   "新建发表",
					Content: "新建发表",
					Author: domain.Author{
						Id:   123,
						Name: "123",
					},
				}).Return(int64(0), errors.New("制作库失败"))
				return author, reader
			},
			wantId:  1,
			wantErr: errors.New("制作库失败"),
		},
		{
			name: "制作库修改失败",
			art: domain.Article{
				Id:      1,
				Title:   "修改发表",
				Content: "修改发表",
				Author: domain.Author{
					Id:   123,
					Name: "123",
				},
			},
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := articlemock.NewMockArticleAuthorRepository(ctrl)
				reader := articlemock.NewMockArticleReaderRepository(ctrl)
				ctx := context.Background()
				author.EXPECT().Update(ctx, domain.Article{
					Id:      1,
					Title:   "修改发表",
					Content: "修改发表",
					Author: domain.Author{
						Id:   123,
						Name: "123",
					},
				}).Return(int64(0), errors.New("制作库失败"))
				return author, reader
			},
			wantId:  1,
			wantErr: errors.New("制作库失败"),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tt.mock(ctrl)
			h := NewArticleServiceV1(author, reader, logger.NewNoOpLogger())
			id, err := h.Publish(context.Background(), tt.art)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, id, tt.wantId)
			require.NoError(t, err)
		})
	}
}
