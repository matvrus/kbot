APP=$(shell basename $(shell git remote get-url origin))
REGISTRY := matvrus
REGISTRYDOC := matvrus
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS=linux #linux darwin windows
TARGETARCH=arm64 #amd64 arm64

format:
	gofmt -s -w ./

lint:
	golint

test:
	go test -v 

get:
	go get

build: format get
	CGO_ENABLED=${CGO_ENABLED} GOOS=${OS} GOARCH=${ARCH} go build -v -o kbot -ldflags "-X="github.com/matvrus/kbot/cmd.appVersion=${VERSION}

image:
	docker build . -t ${REGISTRY}:${VERSION}-${OS}-${ARCH} --build-arg CGO_ENABLED=${CGO_ENABLED} --build-arg ARCH=${ARCH} --build-arg OS=${OS}

push:
	docker push ${REGISTRY}:${VERSION}-${OS}-${ARCH}

pushdoc:
	docker push ${REGISTRYDOC}/${APP}:${VERSION}-${OS}-${ARCH}

imagedoc:
	docker build . -t ${REGISTRYDOC}/${APP}:${VERSION}-${OS}-${ARCH} --build-arg CGO_ENABLED=${CGO_ENABLED} --build-arg TARGETARCH=${ARCH} --build-arg TARGETOS=${TOS}

clean:
	rm -rf kbot
	docker rmi ${REGISTRY}:${VERSION}-${OS}-${ARCH}