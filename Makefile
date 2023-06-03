APP := $(shell basename $(shell git remote get-url origin))
REGISTRY := matvrus
#REGISTRYDOC := matvrus
VERSION := $(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS := linux
TARGETARCH := $(shell dpkg --print-architecture) #amd64 #arm64 

format:
	gofmt -s -w ./

lint:
	golint

test:
	go test -v 

get:
	go get

build: format get
	CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v -o kbot -ldflags "-X=github.com/matvrus/kbot/cmd.appVersion=$(VERSION)"

image:
	docker build . -t ${REGISTRY}:${VERSION}-$(TARGETOS)-${TARGETARCH} 

push:
	docker push ${REGISTRY}:${VERSION}-$(TARGETOS)-${TARGETARCH}
	
pushdoc:
	docker push ${REGISTRYDOC}/${APP}:${VERSION}-${OS}-${ARCH}

imagedoc:
	docker build . -t ${REGISTRYDOC}/${APP}:${VERSION}-${OS}-${ARCH} --build-arg CGO_ENABLED=${CGO_ENABLED} --build-arg TARGETARCH=${ARCH} --build-arg TARGETOS=${TOS}

clean:
	rm -rf kbot
	docker rmi $(REGISTRYDOC)/$(APP):$(VERSION)-$(TARGETOS)-$(TARGETARCH) || true
#comment
