package article

import (
	"context"
	"webookProgram/webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) error
}
