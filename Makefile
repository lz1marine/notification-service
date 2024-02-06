export GO111MODULE := on
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin
REPO := lz1marine/notification-service

SERVER_TAG := server-v0.0.7
WORKER_TAG := worker-v0.0.3
GARBAGE_COLLECTOR_TAG := gc-v0.0.1

.PHONY: ci 
ci: build-all lint test

.PHONY: install-requirements
install-requirements:
	go install github.com/swaggo/swag/cmd/swag@v1.16.3
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: deps
deps:
	go mod vendor
	go mod tidy

.PHONY: fmt
fmt: install-requirements
	$(GOBIN)/goimports -l -w ./

.PHONY: lint
lint: fmt install-requirements
	$(GOBIN)/golangci-lint run --timeout 10m

.PHONY: vet
vet:
	go vet ./...

.PHONY: docs
docs: install-requirements
	$(GOBIN)/swag init -dir ./cmd/notification-service,./pkg/handler,./pkg/handler/inter,./api/v1,./vendor/github.com/swaggo/swag/example/celler/httputil -o ./openapi

.PHONY: tests
tests:
	go test ./...

.PHONY: build-all
build-all: build-notification-server build-notification-worker

.PHONY: build-notification-server
build-notification-server:
	docker build . \
		-f image/notification-server/Dockerfile \
		-t $(REPO):$(SERVER_TAG)
	docker push $(REPO):$(SERVER_TAG)

.PHONY: build-notification-worker
build-notification-worker:
	docker build . \
		-f image/notification-worker/Dockerfile \
		-t $(REPO):$(WORKER_TAG)
	docker push $(REPO):$(WORKER_TAG)

.PHONY: build-notification-gc
build-notification-gc:
	docker build . \
		-f image/notification-worker/Dockerfile \
		-t $(REPO):$(GARBAGE_COLLECTOR_TAG)
	docker push $(REPO):$(GARBAGE_COLLECTOR_TAG)