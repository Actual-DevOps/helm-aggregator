.PHONY: all push_image build_image

APP=helm-aggregator
GROUP=actual-devops
VERSION=0.1.0
DOCKER_REGISTRY=ghcr.io
GOLANG_VERSION=1.23.5

all:
	@echo 'DEFAULT:                                                         '
	@echo '   make build_image                                              '
	@echo '   make build                                                    '
	@echo '   make push_image                                               '

lint:
	golangci-lint run -v

build:
	go mod download
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o helm-aggregator

build_image:
	@echo 'Build Docker'
	docker buildx build --build-arg GOLANG_VERSION=$(GOLANG_VERSION) \
						--platform linux/amd64 \
						-t $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION) .
	docker tag $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION) $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest

push_image:
	docker push $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest
