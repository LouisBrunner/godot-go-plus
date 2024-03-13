all:
.PHONY: all

build: build-macos build-linux build-windows
.PHONY: build

build-macos: build-macos-amd64 build-macos-arm64
.PHONY: build-macos

build-macos-amd64:
	GOOS=darwin GOARCH=amd64 make build-actual
.PHONY: build-macos-amd64

build-macos-arm64:
	GOOS=darwin GOARCH=arm64 make build-actual
.PHONY: build-macos-arm64

build-linux:
	GOOS=linux GOARCH=amd64 make build-actual
.PHONY: build-linux

build-windows:
	GOOS=windows GOARCH=amd64 make build-actual
.PHONY: build-windows

EXTENSION ?= dylib
ifeq (${GOOS},windows)
EXTENSION=dll
else ifeq (${GOOS},linux)
EXTENSION=so
endif

build-actual:
	CGO_ENABLED=1 go build -o libexample-${GOOS}-${GOARCH}.${EXTENSION} -buildmode=c-shared .
.PHONY: build-actual

test-macos:
	/Applications/Godot.app/Contents/MacOS/Godot project.godot --headless -v
.PHONY: test-macos

test-macos-arm64: build-macos-arm64 test-macos
.PHONY: test-macos-arm64
