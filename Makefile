BINARY=issues2markdown-api
MAIN_PACKAGE=cmd/${BINARY}/main.go
PACKAGES = $(shell go list ./...)
VERSION=`cat VERSION`
BUILD=`git symbolic-ref HEAD 2> /dev/null | cut -b 12-`-`git log --pretty=format:%h -1`
DIST_FOLDER=dist
DIST_INCLUDE_FILES=README.md ROADMAP.md LICENSE VERSION

# Setup the -ldflags option for go build here, interpolate the variable
# values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

# Build & Install

install:	## Build and install package on your system
	go install $(LDFLAGS) -v $(PACKAGES)

.PHONY: version
version:	## Show version information
	@echo $(VERSION)-$(BUILD)

# Testing

.PHONY: test
test:		## Execute package tests 
	go test -v $(PACKAGES)

.PHONY: cover-profile
cover-profile:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	rm -rf coverage.out

.PHONY: cover
cover: cover-profile
cover: 		## Generate test coverage data
	go tool cover -func=coverage-all.out

.PHONY: cover-html
cover-html: cover-profile
	go tool cover -html=coverage-all.out

.PHONY: coveralls
coveralls:
	goveralls -service circle-ci -repotoken JoNRiR2QQIzmLZVycbCak5XxegtbMG6Ap

# Lint

lint:		## Lint source code
	gometalinter --disable-all --enable=errcheck --enable=vet --enable=vetshadow

# Dependencies

deps:		## Install package run dependencies
	go get -t -d -u github.com/spf13/cobra/cobra
	go get -t -d -u github.com/issues2markdown/issues2markdown
	go get -u github.com/gorilla/mux
	go get -u golang.org/x/oauth2

dev-deps: deps
dev-deps:	## Install package dev and run dependencies
	go get -t -u github.com/mattn/goveralls
	go get -t -u github.com/inconshreveable/mousetrap
	go get -t -u github.com/alecthomas/gometalinter
	gometalinter --install

# Cleaning up

.PHONY: clean
clean:		## Delete generated development environment
	go clean
	rm -rf ${BINARY}
	rm -rf ${BINARY}.exe
	rm -rf coverage-all.out

# Docs

godoc-serve:	## Serve documentation (godoc format) for this package at port HTTP 9090
	godoc -http=":9090"

# Distribution

.PHONY: dist
dist: clean dist-prepare dist-darwin dist-linux dist-windows	
dist: 		## Generate distribution packages

dist-prepare:
	mkdir -p dist

dist-darwin:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY} ${MAIN_PACKAGE}
	zip ${DIST_FOLDER}/${BINARY}-${VERSION}-darwin-amd64.zip ${BINARY} ${DIST_INCLUDE_FILES}
	GOOS=darwin GOARCH=386 go build ${LDFLAGS} -o ${BINARY} ${MAIN_PACKAGE}
	zip ${DIST_FOLDER}/${BINARY}-${VERSION}-darwin-386.zip ${BINARY} ${DIST_INCLUDE_FILES}
	rm -rf ${BINARY}

dist-linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY} ${MAIN_PACKAGE}
	zip ${DIST_FOLDER}/${BINARY}-${VERSION}-linux-amd64.zip ${BINARY} ${DIST_INCLUDE_FILES}
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o ${BINARY} ${MAIN_PACKAGE}
	zip ${DIST_FOLDER}/${BINARY}-${VERSION}-linux-386.zip ${BINARY} ${DIST_INCLUDE_FILES}
	rm -rf ${BINARY}

dist-windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}.exe ${MAIN_PACKAGE}
	zip ${DIST_FOLDER}/${BINARY}-${VERSION}-windows-amd64.zip ${BINARY}.exe ${DIST_INCLUDE_FILES}
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ${BINARY}.exe ${MAIN_PACKAGE}
	zip ${DIST_FOLDER}/${BINARY}-${VERSION}-windows-386.zip ${BINARY}.exe ${DIST_INCLUDE_FILES}
	rm -rf ${BINARY}.exe

dist-clean: clean 	
	rm -rf ${DIST_FOLDER}

include Makefile.help.mk