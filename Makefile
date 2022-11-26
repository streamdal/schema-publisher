VERSION ?= $(shell git rev-parse --short HEAD)
SERVICE = publish

GO = CGO_ENABLED=$(CGO_ENABLED) go
CGO_ENABLED ?= 0
GO_BUILD_FLAGS = -ldflags "-X main.version=${VERSION}"

# Pattern #1 example: "example : description = Description for example target"
# Pattern #2 example: "### Example separator text
help: HELP_SCRIPT = \
	if (/^([a-zA-Z0-9-\.\/]+).*?: description\s*=\s*(.+)/) { \
		printf "\033[34m%-40s\033[0m %s\n", $$1, $$2 \
	} elsif(/^\#\#\#\s*(.+)/) { \
		printf "\033[33m>> %s\033[0m\n", $$1 \
	}

.PHONY: help
help:
	@perl -ne '$(HELP_SCRIPT)' $(MAKEFILE_LIST)

### Build

.PHONY: build
build: description = Build all
build: clean build/linux build/darwin/amd64 build/darwin/arm64

.PHONY: build/linux
build/linux: description = Build linux
build/linux: clean
	GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o ./$(SERVICE)-linux $(SERVICE)/*.go

.PHONY: build/darwin/amd64
build/darwin/amd64: description = Build darwin for amd64
build/darwin/amd64: clean
	GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o ./$(SERVICE)-darwin-amd64 $(SERVICE)/*.go

.PHONY: build/darwin/arm64
build/darwin/arm64: description = Build darwin for arm64
build/darwin/arm64: clean
	GOOS=darwin GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -o ./$(SERVICE)-darwin-arm64 $(SERVICE)/*.go

.PHONY: clean
clean: description = Remove existing build artifacts
clean:
	$(RM) ./$(SERVICE)-*
