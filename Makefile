#
#  Makefile for Go
#
SHELL=/usr/bin/env bash

GOCMD = go
GOPATH := ${HOME}/go:${GOPATH}
PATH := ${HOME}/bin:${HOME}/src/go/bin:${PATH}
export GOCMD GOPATH

VERSION=$(shell git describe --tags --always)

.PHONY: list

default: build

build:
	${GOCMD} build -ldflags="-X main.Version ${VERSION}" -o bin/wmc-${VERSION}

clean:
	rm bin/*

cross:
	# Not supported by tamr/serial yet
	#echo "Building darwin-amd64..."
	#GOOS="darwin" GOARCH="amd64" go build -ldflags="-X main.Version ${VERSION}" -o bin/wmc-mac-amd64-${VERSION}

	echo "Building windows-386..."
	GOOS="windows" GOARCH="386" go build -ldflags="-X main.Version ${VERSION}" -o bin/wmc-32-${VERSION}.exe

	# Not supported by tamr/serial yet
	#echo "Building freebsd-386..."
	#GOOS="freebsd" GOARCH="386" go build -ldflags="-X main.Version ${VERSION}" -o bin/wmc-freebsd-386-${VERSION}

	echo "Building linux-arm..."
	GOOS="linux" GOARCH="arm" go  build -ldflags="-X main.Version ${VERSION}" -o bin/wmc-linux-arm-${VERSION}

	echo "Building linux-386..."
	GOOS="linux" GOARCH="386" go build -ldflags="-X main.Version ${VERSION}" -o bin/wmc-linux-386-${VERSION}

deploy: cross
	echo "Uploading..."
	ssh thehelpfulhacker.net "mkdir -p ~/media.thehelpfulhacker.net/wmc/${VERSION}"
	scp -v bin/*${VERSION}* thehelpfulhacker.net:~/media.thehelpfulhacker.net/wmc/${VERSION}/
