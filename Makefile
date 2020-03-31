GOPATH ?= $(HOME)/go
LINTER=$(GOPATH)/bin/golangci-lint
GO111MODULE=on

.PHONY: docker_image deps dev_db_up dev_api_up test githooks

$(LINTER):
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint

deps: githooks
	GO111MODULE=on go mod graph
	GO111MODULE=on go mod vendor

lint: $(LINTER)
	GO111MODULE=on $(GOPATH)/bin/golangci-lint run

http_service:
	go build -o http_service .

dev_up:
	docker-compose -f deployments/docker-compose.yaml build --no-cache ideasymbols-db
	docker-compose -f deployments/docker-compose.yaml up -d ideasymbols-db
	docker-compose -f deployments/docker-compose.yaml build --no-cache ideasymbols-http
	docker-compose -f deployments/docker-compose.yaml up -d ideasymbols-http

dev_down:
	docker-compose -f deployments/docker-compose.yaml down

.git/hooks:
	mkdir -p .git/hooks

.git/hooks/pre-push: githooks/pre-push
	cp githooks/pre-push .git/hooks/pre-push

.git/hooks/pre-commit: githooks/pre-commit
	cp githooks/pre-commit .git/hooks/pre-commit

githooks: .git/hooks .git/hooks/pre-push .git/hooks/pre-commit $(LINTER)
	chmod +x .git/hooks/*
