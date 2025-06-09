package article

import (
	"context"
	"webookProgram/webook/internal/domain"
)

type ArticleAuthorRepository interface {
	Update(ctx context.Context, art domain.Article) error
	Create(ctx context.Context, art domain.Article) (int64, error)
}
