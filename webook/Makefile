.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f vincentsun/webook:v0.0.1
	@docker build -t vincentsun/webook:v0.0.1 .