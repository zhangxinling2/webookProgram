package article

import (
	"context"
	"webookProgram/webook/internal/domain"
	"webookProgram/webook/internal/repository/article"
	"webookProgram/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}
type articleServiceV1 struct {
	authorRepo article.ArticleAuthorRepository
	readerRepo article.ArticleReaderRepository
	l          logger.LoggerV1
}

func NewArticleServiceV1(authorRepo article.ArticleAuthorRepository, readerRepo article.ArticleReaderRepository, l logger.LoggerV1) ArticleService {
	return &articleServiceV1{authorRepo, readerRepo, l}
}
func (a *articleServiceV1) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.authorRepo.Update(ctx, art)
		if err != nil {
			return 0, err
		}
		return art.Id, nil
	}
	return a.authorRepo.Create(ctx, art)
}

func (a *articleServiceV1) Publish(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  int64
		err error
	)
	id = art.Id
	if art.Id == 0 {
		id, err = a.authorRepo.Create(ctx, art)
	} else {
		err = a.authorRepo.Update(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			a.l.Error("部分失败:保存到线上数据库失败", logger.Field{Key: "art_id", Value: art.Id}, logger.Error(err))
		} else {
			break
		}
	}

	if err != nil {
		a.l.Error("保存到线上数据库重试全部失败", logger.Field{Key: "art_id", Value: art.Id}, logger.Error(err))
		return 0, err
	}
	return id, nil
}

type articleService struct {
	repo article.ArticleRepository
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}
func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		if err != nil {
			return 0, err
		}
		return art.Id, nil
	}
	return a.repo.Create(ctx, art)
}
func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}
