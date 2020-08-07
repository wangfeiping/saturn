PACKAGES=$(shell go list ./... | grep -v '/simulation')
DIR_BUILD=./build

NAME := Saturn
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

# TODO: Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(NAME) \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=saturnd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=saturncli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: test
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/saturnd
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/saturncli

build: test
		go build -mod=readonly $(BUILD_FLAGS) -o ${DIR_BUILD}/saturnd   ./cmd/saturnd
		go build -mod=readonly $(BUILD_FLAGS) -o ${DIR_BUILD}/saturncli ./cmd/saturncli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test: go.sum
		@if [ -d ${DIR_BUILD} ]; then exit 0; else mkdir ${DIR_BUILD}; fi
		@go test -mod=readonly -coverprofile=build/covprofile $(PACKAGES)
		@go tool cover -html=build/covprofile -o build/coverage.html

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
