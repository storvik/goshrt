# Version=`git describe --tags`  # git tag 1.0.1  # require tag tagged before
GitVersion     := 0.0.1
GitCommit      := `git describe --abbrev=8 --dirty --always`
BuildTime      := `date +%Y-%m-%d\ %H:%M`
BuildGoVersion := `go version`
BUILD_DIR   := build

LDFLAGS := -X 'github.com/storvik/goshrt/version.GitVersion=${GitVersion}'
LDFLAGS += -X 'github.com/storvik/goshrt/version.GitCommit=${GitCommit}'
LDFLAGS += -X 'github.com/storvik/goshrt/version.BuildTime=${BuildTime}'
LDFLAGS += -X 'github.com/storvik/goshrt/version.BuildGoVersion=${BuildGoVersion}'

.PHONY: \
	help \
	all \
	build \
	clean \
	lint \
	test \
	version \
	vet

.DEFAULT_GOAL := build

clean: ## Clean project, run go clean and delete build/
	rm -rf $(BUILD_DIR)

test: ## Run tests using go test
	go test -v ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

vet: ## Reports suspicious constructs etc.
	go vet -v ./...

errors: ## Check for unchecked errors
	errcheck -ignoretests -blank ./...

lint: ## Run golint linter
	golint ./...

version: ## Print project version string
	@echo "$(GitVersion)"
	@echo "$(GitCommit)"
	@echo "$(BuildTime)"
	@echo "$(BuildGoVersion)"

build: ## Build application for current system
	mkdir -p $(BUILD_DIR)/goshrt
	mkdir -p $(BUILD_DIR)/goshrtc
	go build -v \
		-ldflags "-w -s $(LDFLAGS)" \
		-o $(BUILD_DIR)/goshrt/goshrt \
		cmd/goshrt/*.go
	go build -v \
		-ldflags "-w -s $(LDFLAGS)" \
		-o $(BUILD_DIR)/goshrtc/goshrtc \
		cmd/goshrtc/*.go

help: ## Display Makefile help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
