package article

import (
	"context"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/dao/article"
)

type ArticleAuthorRepository interface {
	Update(ctx context.Context, art domain.Article) error
	Create(ctx context.Context, art domain.Article) (int64, error)
}
type articleAuthorRepository struct {
	dao article.AuthorDAO
}

func NewArticleAuthorRepository(dao article.AuthorDAO) ArticleAuthorRepository {
	return &articleAuthorRepository{dao: dao}
}
func (a *articleAuthorRepository) Update(ctx context.Context, art domain.Article) error {
	return a.dao.UpdateById(ctx, a.domainToEntity(art))
}

func (a *articleAuthorRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Insert(ctx, a.domainToEntity(art))
}
func (c *articleAuthorRepository) entityToDomain(art article.Article) domain.Article {
	return domain.Article{
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
			//Name: "",
		},
	}
}
func (c *articleAuthorRepository) domainToEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
