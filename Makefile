PROJECT_NAME		:= bse_shares
HOST_DIRECTORY		:= output
GIT_TAG				:= $(shell git describe --dirty --tags --always)
GIT_COMMIT			:= $(shell git rev-parse --short HEAD)
LDFLAGS				:= -X "main.gitTag=$(GIT_TAG)" -X "main.gitCommit=$(GIT_COMMIT)" -extldflags "-static" -s -w

FIRST_GOPATH			:= $(firstword $(subst :, ,$(shell go env GOPATH)))
GOLANGCI_LINT_BIN		:= $(FIRST_GOPATH)/bin/golangci-lint

.PHONY: all
all: vendor build

.PHONY: clean
clean:
	git clean -Xfd .

#######################################
# builds
#######################################

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: build-all
build-all:
	GOOS=linux   GOARCH=${GOARCH} CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o '$(PROJECT_NAME)' .
	GOOS=darwin  GOARCH=${GOARCH} CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o '$(PROJECT_NAME).darwin' .
	GOOS=windows GOARCH=${GOARCH} CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o '$(PROJECT_NAME).exe' .

.PHONY: build
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(PROJECT_NAME) .

.PHONY: image
image: image
	docker build -t $(PROJECT_NAME):$(GIT_TAG) .

.PHONY: build-push-development
build-push-development:
	docker buildx create --use
	docker buildx build -t vibin/$(PROJECT_NAME):development --platform linux/amd64,linux/arm,linux/arm64 --push .

#######################################
# quality checks
#######################################

.PHONY: check
check: vendor lint test

.PHONY: test
test:
	time go test ./...

.PHONY: lint
lint: $(GOLANGCI_LINT_BIN)
	time $(GOLANGCI_LINT_BIN) run --verbose --print-resources-usage

$(GOLANGCI_LINT_BIN):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(FIRST_GOPATH)/bin