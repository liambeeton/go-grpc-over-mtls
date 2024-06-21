REPO_NAME = go-grpc-over-mtls
REPO = github.com/liambeeton/go-grpc-over-mtls
CLIENT_NAME = ${REPO_NAME}-client
CLIENT_IMAGE = github.com/liambeeton/go-grpc-over-mtls/client
FQDN = localhost
BUILD_TIME = $(shell date -u +"%FT%TZ")
GOLANGCI_LINT = v1.59.1-alpine
PRETTIER = 3.1.0
PROTOC_VERSION ?= 24.0

.DEFAULT_GOAL := build

.PHONY: build
build:
	GOBIN=${PWD}/bin go install \
		-mod=vendor \
		-ldflags="-X main.buildTime=${BUILD_TIME}" \
		./cmd/...

.PHONY: linux
linux:
	CGO_ENABLED=0 \
	GOOS=linux \
	GO111MODULE=auto \
	GOBIN=${CURDIR} \
	go install \
		-mod=vendor \
		-ldflags="-X main.buildTime=${BUILD_TIME}" \
		-trimpath \
		./cmd/...

.PHONY: test
test:
	go test ./... -cover -coverprofile coverage.out -race -mod=vendor -v

.PHONY: run
run:
	go run \
		-mod=vendor \
		-ldflags="-X main.buildTime=${BUILD_TIME}" \
		./cmd/server \
		--host=${FQDN} \
		--port=8443 \
		--key-file=./certs-rsa/server.key \
		--cert-file=./certs-rsa/server.crt \
		--ca-file=./certs-rsa/ca.crt

.PHONY: build-docker
build-docker:
	docker build -t ${REPO} .

.PHONY: run-docker
run-docker:
	docker run -p 8443:8443 ${REPO}

.PHONY: run-client
run-client:
	go run \
		-mod=vendor \
		-ldflags="-X main.buildTime=${BUILD_TIME}" \
		cmd/client/*.go \
		--host=${FQDN} \
		--port=8443 \
		--key-file=./certs-rsa/client.key \
		--cert-file=./certs-rsa/client.crt \
		--ca-file=./certs-rsa/ca.crt

.PHONY: build-docker-client
build-docker-client:
	docker build -f client.Dockerfile -t ${CLIENT_IMAGE}:local .

.PHONY: run-docker-client
run-docker-client:
	docker run ${REPO}

.PHONY: lint
lint:
	docker run \
		--rm \
		-t \
		-w /app \
		-v ${PWD}:/app \
		golangci/golangci-lint:${GOLANGCI_LINT} \
		golangci-lint run \
			-v \
			-c .golangci.yml

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: upgrade-vendor
upgrade-vendor:
	go get -u ./...

.PHONY: gofmt
gofmt:
	go fmt ./...

.PHONY: markdownlint
markdownlint:
	docker run --rm -v "${PWD}:/work" ghcr.io/tmknom/dockerfiles/markdownlint:0.37.0 -- .

.PHONY: prettier-check
prettier-check:
	docker run --rm -v "${PWD}:/work" ghcr.io/tmknom/dockerfiles/prettier:${PRETTIER} --check .

.PHONY: prettier-fix
prettier-fix:
	docker run --rm -v "${PWD}:/work" ghcr.io/tmknom/dockerfiles/prettier:${PRETTIER} --write .

.PHONY: proto-clean
proto-clean:
	@rm -rf pb/message
	@rm -rf pb/service

.PHONY: proto-compile
proto-compile:
	PLATFORM=$(shell uname -m) PROTOC_VERSION=$(PROTOC_VERSION) docker-compose -f docker/docker-compose.yml run --rm protogen

.PHONY: docker-config
docker-config:
	PLATFORM=$(shell uname -m) PROTOC_VERSION=$(PROTOC_VERSION) docker-compose -f docker/docker-compose.yml config
