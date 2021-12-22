PROJECT := minikube
GIT_COMMIT := $(shell git rev-parse HEAD)
VERSION := latest
APP_NAME := k8s-go-app

build:
	go build -ldflags "-X $(PROJECT)/version.Version=1.0.0 \
	-X $(PROJECT)/version.Commit=$(GIT_COMMIT) " -o ./bin/app main.go

build_container:
	docker build --build-arg=GIT_COMMIT=$(GIT_COMMIT) \
                 --build-arg=VERSION=$(VERSION) \
                 --build-arg=PROJECT=$(PROJECT) \
                 -t docker.io/ppdraga/simple_go_app .
