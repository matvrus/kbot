APP := $(shell basename $(shell git remote get-url origin))
REGISTRY := matvrus
REGISTRYDOC := matvrus
VERSION := $(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS := linux
TARGETARCH := $(shell dpkg --print-architecture)

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
	docker build . -t $(REGISTRY)/$(APP):$(VERSION)-$(TARGETARCH) --build-arg TARGETARCH=$(TARGETARCH)

push:
	docker push $(REGISTRY)/$(APP):$(VERSION)-$(TARGETARCH)
	
pushdoc:
	docker push $(REGISTRYDOC)/$(APP):$(VERSION)-$(TARGETOS)-$(TARGETARCH)

imagedoc:
	docker build . -t $(REGISTRYDOC)/$(APP):$(VERSION)-$(TARGETOS)-$(TARGETARCH) --build-arg CGO_ENABLED= --build-arg TARGETARCH=$(TARGETARCH) --build-arg TARGETOS=$(TARGETOS)

clean:
	rm -rf kbot
	docker rmi $(REGISTRY)/$(APP):$(VERSION)-$(TARGETARCH)
