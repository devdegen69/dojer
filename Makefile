.PHONY: build-ui
build-ui:
	cd ui && npm run build

.PHONY: build-go
build-go:
	go build -ldflags '-linkmode external -extldflags "-static"' -tags 'release sqlite_math_functions'

.PHONY: build-all
build-all: build-ui build-go

.DEFAULT_GOAL := build-all
