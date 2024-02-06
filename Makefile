APP=$(shell basename -s .git $(shell git remote get-url origin))
REGISTRY ?= ghcr.io/yuandrk
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS ?= linux
TARGETARCH ?= amd64

format:
	gofmt -s -w ./

get:
	go get

lint:
	golint

test:
	go test -v

build: format get
	CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v -o sbot -ldflags "-X="github.com/yuandrk/sbot/cmd.appVersion=${VERSION}

linux: format get
	CGO_ENABLED=0 GOOS=linux GOARCH=$(TARGETARCH) go build -v -o sbot -ldflags "-X="github.com/yuandrk/sbot/cmd.appVersion=${VERSION}
	docker build --build-arg name=linux -t ${REGISTRY}/${APP}:${VERSION}-linux-${TARGETARCH} .

windows: format get
	CGO_ENABLED=0 GOOS=windows GOARCH=$(TARGETARCH) go build -v -o sbot -ldflags "-X="github.com/yuandrk/sbot/cmd.appVersion=${VERSION}
	docker build --build-arg name=windows -t ${REGISTRY}/${APP}:${VERSION}-windows-$(TARGETARCH) .

darwin:format get
	CGO_ENABLED=0 GOOS=darwin GOARCH=$(TARGETARCH) go build -v -o sbot -ldflags "-X="github.com/yuandrk/sbot/cmd.appVersion=${VERSION}
	docker build --build-arg name=darwin -t ${REGISTRY}/${APP}:${VERSION}-darwin-$(TARGETARCH) .

arm: format get
	CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=arm go build -v -o sbot -ldflags "-X="github.com/yuandrk/sbot/cmd.appVersion=${VERSION}
	docker build --build-arg TARGETARCH=arm -t ${REGISTRY}/${APP}:${VERSION}-$(TARGETOS)-arm .

image: 
	docker build . --build-arg TARGETARCH=$(TARGETARCH) -t ${REGISTRY}/${APP}:${VERSION}-$(TARGETOS)-$(TARGETARCH) 

push:
	docker push ${REGISTRY}/${APP}:${VERSION}-$(TARGETOS)-$(TARGETARCH)

clean:
	@rm -rf sbot; \
	IMG1=$$(docker images -q | head -n 1); \
	if [ -n "$${IMG1}" ]; then  docker rmi -f $${IMG1}; else printf "$RImage not found$D\n"; fi


