.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmock -destination=webook/internal/service/mocks/user.mock.go
	@mockgen -source=webook/internal/service/code.go -package=svcmock -destination=webook/internal/service/mocks/code.mock.go
	@mockgen -source=webook/internal/service/article/article.go -package=artmock -destination=webook/internal/service/mocks/article/article.mock.go
	@mockgen -source=webook/internal/repository/user.go -package=svcmock -destination=webook/internal/repository/mocks/user.mock.go
	@mockgen -source=webook/internal/repository/code.go -package=svcmock -destination=webook/internal/repository/mocks/code.mock.go
	@mockgen -source=webook/internal/repository/cache/user.go -package=cachemock -destination=webook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=webook/internal/repository/dao/user.go -package=daomock -destination=webook/internal/repository/dao/mocks/user.mock.go
    @mockgen -source=webook/internal/repository/article/author.go -package=articlemock -destination=webook/internal/repository/article/mocks/article/author.mock.go
    @mockgen -source=webook/internal/repository/article/reader.go -package=articlemock -destination=webook/internal/repository/article/mocks/article/reader.mock.go
    @mockgen -source=webook/internal/repository/dao/article/author.go -package=articledaomock -destination=webook/internal/repository/dao/article/mocks/article/author.mock.go
    @mockgen -source=webook/internal/repository/dao/article/reader.go -package=articledaomock -destination=webook/internal/repository/dao/article/mocks/article/reader.mock.go


	@mockgen -package=redismock -destination=webook/internal/repository/cache/redismocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable
	@mockgen -source=webook/internal/repository/code.go -package=svcmock -destination=webook/internal/repository/mocks/code.mock.go
	@go mod tidy