package article

import (
	"context"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/dao/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	FindById(ctx context.Context, articleId int64) (domain.Article, error)
	Sync(ctx context.Context, art domain.Article) (int64, error)
}
type CacheArticleRepository struct {
	//操作单一的库
	dao article.ArticleDAO
	//操作不同库
	authorDAO article.AuthorDAO
	readerDAO article.ReaderDAO
	//cache ArticleCache
}

func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.domainToEntity(art))
}
func (c *CacheArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = c.authorDAO.Insert(ctx, c.domainToEntity(art))
	} else {
		err = c.authorDAO.UpdateById(ctx, c.domainToEntity(art))
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	err = c.readerDAO.UpsertArticle(ctx, c.domainToEntity(art))
	return id, err
}
func (c *CacheArticleRepository) FindById(ctx context.Context, articleId int64) (domain.Article, error) {
	ae, err := c.dao.FindByID(ctx, articleId)
	if err != nil {
		return domain.Article{}, err
	}
	return c.entityToDomain(ae), err
}

func NewArticleRepository(dao article.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{dao: dao}
}

func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, c.domainToEntity(art))
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.domainToEntity(art))
}
func (c *CacheArticleRepository) entityToDomain(art article.Article) domain.Article {
	return domain.Article{
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
			//Name: "",
		},
	}
}
func (c *CacheArticleRepository) domainToEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
