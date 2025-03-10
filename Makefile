.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmock -destination=webook/internal/service/mocks/user.mock.go
	@mockgen -source=webook/internal/service/code.go -package=svcmock -destination=webook/internal/service/mocks/code.mock.go
	@go mod tidy