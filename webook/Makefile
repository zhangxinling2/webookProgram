.PHONY: docker
docker:
	@rm webook || true
	@go build -o webook .
	@docker rmi -f zhangxinling2/webook:v0.0.1
	@docker build -t zhangxinling2/webook:v0.0.1 .