VERSION := 0.0.1
COMMIT :=$(shell git rev-parse HEAD)
BRANCH :=$(shell git rev-parse --abbrev-ref HEAD)
BIN_DIR := $(shell pwd)/build
CURR_DIR :=$(shell pwd)


# Setup the -ldflags option to pass vars defined here to app vars
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH} -extldflags \"-static\""
GOPATH= ${CURR_DIR}

PKGS := $(shell go list ./... | grep -v /vendor)

PLATFORMS := windows linux darwin
os = $(word 1, $@)

build:
	go get ./...
	go build sfill
	CGO_ENABLED=0 GOOS=linux go build -o template -a -ldflags '-extldflags "-static"' src/sfill/scripts/template/template.go
	CGO_ENABLED=0 GOOS=linux go build -o template -a -ldflags '-extldflags "-static"' src/sfill/scripts/loader/loader.go


.PHONY: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build ${LDFLAGS} -o $(BIN_DIR)/$(BINARY)-$(VERSION)-$(os)-amd64

test:
	go test github.com/vorsprung/parameter-ssm-template/sfill

