.PHONY: all push_image build_image

APP=helm-aggregator
GROUP=actual-devops
VERSION=$(shell cat version)
DOCKER_REGISTRY=ghcr.io
GOLANG_VERSION=1.23.5

BUILD_CMD='GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o helm-aggregator'

all:
	@echo 'DEFAULT:                                                         '
	@echo '   make build_image                                              '
	@echo '   make build                                                    '
	@echo '   make push_image                                               '

lint:
	golangci-lint run -v

build:
	go mod download
	$(shell echo $(BUILD_CMD))

build_image:
	@echo 'Build Docker'
	docker buildx build --build-arg GOLANG_VERSION=$(GOLANG_VERSION) \
						--build-arg BUILD_CMD=$(BUILD_CMD) \
						--platform linux/amd64 \
						-t $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION) .
	docker tag $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION) $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest

push_image:
	docker push $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest
