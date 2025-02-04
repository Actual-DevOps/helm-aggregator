.PHONY: all build push build_app

APP=helm-aggregator
GROUP=Actual-DevOps
VERSION=0.1.0
DOCKER_REGISTRY=ghcr.io


all:
	@echo 'DEFAULT:                                                               '
	@echo '   make build_app                                                      '
	@echo '   make build                                                    '
	@echo '   make push                                                           '

build_app:
	golangci-lint run -v
	go mod vendor
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o helm-aggregator

build:
	@echo 'Build Docker'
	docker buildx build --platform linux/amd64 -t $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION) .
	docker tag $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION) $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest

push:
	docker push $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest
