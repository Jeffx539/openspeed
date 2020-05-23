VERSION   ?= $(shell git describe --tags)
REVISION  ?= $(shell git rev-parse HEAD)
BRANCH    ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILDUSER ?= $(shell id -un)
BUILDTIME ?= $(shell date '+%Y%m%d-%H:%M:%S')

.PHONY: vue-build

build: vue-build
	go build -o bin/openspeed openspeed.go
	mv ${PWD}/web/dist ${PWD}/bin/static

clean:
	rm -rf ./bin

vue-build:
	cd ${PWD}/web; yarn build
	