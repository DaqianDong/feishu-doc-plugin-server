.PHONY: build backend

PROJECT_ROOT=$(shell pwd)
GOBIN_PATH=$(PROJECT_ROOT)/bin

build: switch_linux backend switch_drawin

build_img:build
	cd docker && sh build_images.sh ${version}

switch_linux:
	# 指定go build的编译目标为linux darwin
	go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64

switch_drawin:
	# 还原
	go env -w CGO_ENABLED=1 GOOS=darwin GOARCH=amd64

backend:
	go build -o $(GOBIN_PATH)/server *.go