package article

import (
	"context"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/dao/article"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) error
}
type articleReaderRepository struct {
	dao article.ReaderDAO
}

func NewArticleReaderRepository(dao article.ReaderDAO) ArticleReaderRepository {
	return &articleReaderRepository{dao: dao}
}
func (a *articleReaderRepository) Save(ctx context.Context, art domain.Article) error {
	return a.dao.UpsertArticle(ctx, a.domainToEntity(art))
}

func (a *articleReaderRepository) entityToDomain(art article.Article) domain.Article {
	return domain.Article{
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
			//Name: "",
		},
	}
}
func (a *articleReaderRepository) domainToEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
