MAKEFILE_PATH := $(abspath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
PATH := $(MAKEFILE_PATH):$(PATH)

export GOBIN := $(MAKEFILE_PATH)/bin

all: clean format build test

clean:
	@echo clean
	@go clean

build:
	@echo build
	@go build -o $(GOBIN)/news-api

TEST_PKGS = $(shell go list ./... | grep -v /test)

.PHONY: test
test:
	@echo test
	@go test -count=1 -v $(TEST_PKGS)

lint: install-tools
	@echo lint
	@$(GOBIN)/golangci-lint run

format:
	@echo format
	@go fmt ./...

tidy:
	@echo tidy
	@go mod tidy

IMAGE = news-api

docker:
	@echo docker
	@docker build -t $(IMAGE) -f Dockerfile .

docker-run:
	@echo docker-run
	@docker run --rm $(IMAGE)

install-tools:
	@echo install tools from tools.go
	@cat tools.go | grep _ | awk -F '"' '{print $$2}' | xargs -tI % go install %
