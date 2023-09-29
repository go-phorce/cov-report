include ./.project/go-project.mk

# don't echo execution
.SILENT:

.DEFAULT_GOAL := help

.PHONY: *

default: help

all: tools build test

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54

version:
	gofmt -r '"GIT_VERSION" -> "$(GIT_VERSION)"' version/current.template > version/current.go

build:
	go build -o ${PROJ_ROOT}/bin/cov-report ./cmd/cov-report

