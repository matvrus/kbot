# APP := $(shell basename $(shell git remote get-url origin))
# REGISTRY := matvrus
# #REGISTRYDOC := matvrus
# VERSION := $(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
# TARGETOS := linux
# TARGETARCH := $(shell dpkg --print-architecture) #amd64 #arm64 

APP =$(shell basename $(shell git remote get-url origin))
REGISTRY =matvrus
GO_BUILD =go build -v -o kbot -ldflags "-X"=github.com/matvrus/kbot/cmd.appVersion=
VERSION =$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
# TARGETOS can be: linux darwin windows. For macos don't forget to change TARGETARCH=arm64
TARGETOS =linux
# TARGETARCH can be: amd64 arm64
TARGETARCH ?=amd64
# For windows change 0 to 1
CGO_ENABLED ?=0
TARGETOS_LIST = linux darwin windows
TARGETARCH_LIST = amd64 arm64

.PHONY: all format lint test get build build_all image image_all push push_all clean clean_all_builds clean_all_images

all: format lint test build_all image_all push_all clean_all_builds clean_all_images

format:
	gofmt -s -w ./

lint:
	golint

test:
	go test -v

get:
	go get

build: format get
	CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} ${GO_BUILD}${VERSION}

image:
	docker build . -t ${REGISTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH} --build-arg TARGETOS=${TARGETOS} --build-arg TARGETARCH=${TARGETARCH} --build-arg CGO_ENABLED=${CGO_ENABLED}

build_all:
	for os in $(TARGETOS_LIST); do \
		for arch in $(TARGETARCH_LIST); do \
			if [ "$$os" = "windows" ]; then \
				make build TARGETOS=$$os TARGETARCH=$$arch CGO_ENABLED=1; \
			else \
				make build TARGETOS=$$os TARGETARCH=$$arch; \
			fi; \
   			mkdir -p builds/$${os}_$${arch} && mv kbot builds/$${os}_$${arch}/; \
		done \
	done

image_all:
	for os in $(TARGETOS_LIST); do \
		for arch in $(TARGETARCH_LIST); do \
			if [ "$$os" = "windows" ]; then \
				make image TARGETOS=$$os TARGETARCH=$$arch CGO_ENABLED=1; \
			else \
				make image TARGETOS=$$os TARGETARCH=$$arch; \
			fi; \
		done \
	done

push:
	docker push ${REGISTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH}

push_all:
	for TARGETOS in $(TARGETOS_LIST); do \
		for TARGETARCH in $(TARGETARCH_LIST); do \
			docker push ${REGISTRY}/${APP}:${VERSION}-$$TARGETOS-$$TARGETARCH; \
		done ; \
	done

clean:
	rm -rf kbot
	docker rmi ${REGISTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH}

clean_all_images:
	for TARGETOS in $(TARGETOS_LIST); do \
		for TARGETARCH in $(TARGETARCH_LIST); do \
			docker rmi ${REGISTRY}/${APP}:${VERSION}-$$TARGETOS-$$TARGETARCH; \
		done ; \
	done

clean_all_builds:
	rm -rf ./builds/