# basic info
app := govm
module := github.com/Open-Source-CQUT/govm/cmd/govm
output := $(shell pwd)/bin
# meta info
build_time := $(shell date +"%Y/%m/%dT%H:%M:%SZ%z")
git_version := $(shell git describe --tags --always)
# build info
host_os := $(shell go env GOHOSTOS)
host_arch := $(shell go env GOHOSTARCH)
os := $(host_os)
arch := $(host_arch)

ifeq ($(os), windows)
	exe := .exe
endif


.PHONY: build
build:
	# go lint
	go vet ./...

	# prepare target environment $(os)/$(arch)
	go env -w GOOS=$(os)
	go env -w GOARCH=$(arch)

	# build go module
	go build -a -trimpath \
		-ldflags="-X main.AppName=$(app) -X main.Version=$(git_version) -X main.BuildTime=$(build_time)" \
		-o $(output)/$(app)-$(os)-$(arch)$(exe) \
		$(module)

	# resume host environment $(host_os)/$(host_arch)
	go env -w GOOS=$(host_os)
	go env -w GOARCH=$(host_arch)