VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main. version=$(VERSION)' -X 'main. revision=$(REVISION)'
SRC_DIR := ./
BIN_DIR := bin
BINARY := bin/tldr
WORKFLOW_DIR := "$${HOME}/Library/Application Support/Alfred 3/Alfred.alfredpreferences/workflows/user.workflow.2569C1E1-8114-4B77-9506-52AA966313A9"
ASSETS_DIR := assets
ARTIFACT_DIR := .artifact
ARTIFACT := ${ARTIFACT_DIR}/tldr.alfredworkflow

export GO111MODULE=on

## Build binaries on your environment
build: setup
	CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o ${BINARY} ./${SRC_DIR}

## Setup
setup:
	#installing golint
	@(if ! type golint >/dev/null 2>&1; then go get -u golang.org/x/lint/golint ;fi)
	#installing golangci-lint
	@(if ! type golangci-lint >/dev/null 2>&1; then go get -u github.com/golangci/golangci-lint/cmd/golangci-lint ;fi)
	#installing goimports
	@(if ! type goimports >/dev/null 2>&1; then go get -u golang.org/x/tools/cmd/goimports ;fi)
	#installing ghr
	@(if ! type ghr >/dev/null 2>&1; then go get -u github.com/tcnksm/ghr ;fi)
	#installing make2help
	@(if ! type make2help >/dev/null 2>&1; then go get -u github.com/Songmu/make2help/cmd/make2help ;fi)

## Format source codes
fmt: setup
	goimports -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

## Lint
lint: setup
	golangci-lint run ./...

## Build linux binaries
darwin: setup
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS} -s -w" -o ${BINARY} ./${SRC_DIR}

## Run tests for my project
test: setup
	go test -v ./...

## Install Binary and Assets to Workflow Directory
install: build
	@(cp ${BINARY} ${WORKFLOW_DIR}/)
	@(cp ${ASSETS_DIR}/*  ${WORKFLOW_DIR}/)

## Initialize directory
init:
	@(if [ ! -e ${SRC_DIR} ]; then mkdir ${SRC_DIR}; fi)
	@(if [ ! -e ${BIN_DIR} ]; then mkdir ${BIN_DIR}; fi)
	@(if [ ! -e go.mod ]; then go mod init; fi)

release: darwin
	@(if [ ! -e ${ARTIFACT_DIR} ]; then mkdir ${ARTIFACT_DIR} ; fi)
	@(cp ${BINARY} ${ARTIFACT_DIR})
	@(cp ${ASSETS_DIR}/* ${ARTIFACT_DIR})
	@(zip -j ${ARTIFACT} ${ARTIFACT_DIR}/*)
	@ghr -replace ${VERSION} ${ARTIFACT}

## Clean Binary
clean:
	rm -f ${BIN_DIR}/*
	rm -f ${ARTIFACT_DIR}/*

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: build setup test lint fmt linux init clean help
