package article

import "context"

type ReaderDAO interface {
	UpsertArticle(ctx context.Context, art Article) error
}
