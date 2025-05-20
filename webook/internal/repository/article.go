package repository

import (
	"context"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	FindById(ctx context.Context, articleId int64) (domain.Article, error)
}
type CacheArticleRepository struct {
	dao dao.ArticleDAO
	//cache ArticleCache
}

func (c *CacheArticleRepository) FindById(ctx context.Context, articleId int64) (domain.Article, error) {
	ae, err := c.dao.FindByID(ctx, articleId)
	if err != nil {
		return domain.Article{}, err
	}
	return c.entityToDomain(ae), err
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{dao}
}

func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.Update(ctx, c.domainToEntity(art))
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.domainToEntity(art))
}
func (c *CacheArticleRepository) entityToDomain(art dao.Article) domain.Article {
	return domain.Article{
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
			//Name: "",
		},
	}
}
func (c *CacheArticleRepository) domainToEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
