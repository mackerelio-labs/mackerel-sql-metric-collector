NAME ?= mackerel-sql-metric-collector
REGISTRY_URI ?= $(shell id -u -n)
OUTPUT_DIR ?= bin

GIT_REVISION := $(shell git rev-parse --short HEAD)
TAG := $(GIT_REVISION)

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
BUILD_LDFLAGS := "-X main.revision=$(GIT_REVISION)"

SOURCES = $(shell find . -type f -name '*.go')

all: $(OUTPUT_DIR)/$(NAME)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -v ./...

$(OUTPUT_DIR)/$(NAME): $(OUTPUT_DIR) $(SOURCES)
	go mod tidy
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags=$(BUILD_LDFLAGS) -o $(OUTPUT_DIR)/$(NAME) ./cmd/mackerel-sql-metric-collector/

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

.PHONY: clean
clean: $(OUTPUT_DIR)
	rm -rf $(OUTPUT_DIR)

.PHONY: build-image
build-image:
	docker build -t $(NAME):$(TAG) .

.PHONY: push-image
push-image: build-image
	docker tag $(NAME):$(TAG) $(REGISTRY_URI)/$(NAME):$(TAG)
	docker push $(REGISTRY_URI)/$(NAME):$(TAG)

.PHONY: push-latest-image
push-latest-image: push-image
	docker tag $(REGISTRY_URI)/$(NAME):$(TAG) $(REGISTRY_URI)/$(NAME):latest
	docker push $(REGISTRY_URI)/$(NAME):latest

$(OUTPUT_DIR)/linux/$(NAME): $(SOURCES)
	$(MAKE) $(OUTPUT_DIR)/linux/$(NAME) OS=linux ARCH=amd64 OUTPUT_DIR=$(OUTPUT_DIR)/linux

$(OUTPUT_DIR)/lambda.zip: $(OUTPUT_DIR)/linux/$(NAME)
	zip -j $(OUTPUT_DIR)/lambda.zip $(OUTPUT_DIR)/linux/$(NAME)

.PHONY: lambda-artifact
lambda-artifact: $(OUTPUT_DIR)/lambda.zip

$(OUTPUT_DIR)/bootstrap.zip: $(OUTPUT_DIR)/linux/$(NAME)
	mv $(OUTPUT_DIR)/linux/$(NAME) $(OUTPUT_DIR)/linux/bootstrap
	zip -j $(OUTPUT_DIR)/bootstrap.zip $(OUTPUT_DIR)/linux/bootstrap

.PHONY: lambda-bootstrap-artifact
labmda-boostrap-artifact: $(OUTPUT_DIR)/bootstrap.zip
