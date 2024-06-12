GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ROOT_DIR := app
DIRS := $(wildcard $(ROOT_DIR)/*)
APPS := $(wildcard $(ROOT_DIR)/apps/*)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c find api -type f -name "*.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -type f -name "*.proto")
endif

# Define a dynamic target for each app
.PHONY: $(APPS)
$(APPS):
	@echo "Running $@"
	./bin/$@ &

.PHONY: docker_build
docker_build:
	docker build -t service-orchestrator:latest .

.PHONY: docker_run
docker_run:
	docker run --rm -p 8000:8000 -p 9000:9000 -v ./data/conf service-orchestrator:latest

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go mod tidy

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: wire
wire:
	make config && make api
	@for dir in $(DIRS); do \
		if [ -d $$dir ]; then \
			echo "Running wire in $$dir"; \
			(cd $$dir && wire); \
		fi \
	done

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...
	make proto
	make wire
	

.PHONY: generate
generate:
	go generate ./...

.PHONY: proto
proto:
	make config;
	make api;

.PHONY: all
all: 
	@for dir in $(DIRS); do \
		app_name=$$(basename $$dir); \
		echo "Running $$app_name"; \
		./bin/$$app_name & \
	done
	./bin/core &

.PHONY: execute
execute: $(APPS)
	./bin/core &

.PHONY: compile
compile:
	make init;
	make build;

.PHONY: run
run:
	make build;
	make all;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help