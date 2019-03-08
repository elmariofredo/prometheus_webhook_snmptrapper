APP_NAME  := prometheus_webhook_snmptrapper
DEBUG_POSTFIX := -debugger
NAMESPACE := sysincz
IMAGE := $(NAMESPACE)/$(APP_NAME)
DEBUG_IMAGE :=  $(NAMESPACE)/$(APP_NAME)$(DEBUG_POSTFIX)
REGISTRY := docker.io
ARCH := linux 
VERSION := $(shell git describe --tags 2>/dev/null)
ifeq "$(VERSION)" ""
VERSION := $(shell git rev-parse --short HEAD)
endif
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE=$(shell date +%FT%T%z)
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Branch=$(BRANCH) -X main.BuildDate=$(BUILD_DATE)"

.PHONY: clean

clean:
	rm -rf bin/%/$(APP_NAME)

dep:
	go get -v ./...

build: clean dep
	GOOS=$(ARCH) GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -a -installsuffix cgo -o bin/$(APP_NAME) ./
	GOOS=$(ARCH) GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -a -installsuffix cgo  -o bin/$(APP_NAME)$(DEBUG_POSTFIX) ./trapdebug


docker: build build-trapper-image build-trapper-image-debug

build-trapper-image: build
	docker build . -t $(IMAGE):$(VERSION)

build-trapper-image-debug:  build
	docker build . --file ./trapdebug/Dockerfile -t $(DEBUG_IMAGE):$(VERSION)

docker-push: build-image
	docker push $(IMAGE):$(VERSION)
	docker push $(DEBUG_IMAGE):$(VERSION)


webhook: build-trapper-image
	docker run -p 9099:9099/tcp --network host -v "/tmp/config:/config" --rm --name $(APP_NAME) $(IMAGE):$(VERSION) -debug

debug_webhook: build-trapper-image
	docker run -it -p 9099:9099/tcp -e SNMP_SERVER=localhost --network host --rm --entrypoint "/bin/bash" --name $(APP_NAME) $(IMAGE):$(VERSION) 
