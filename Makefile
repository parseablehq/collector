SHELL := /bin/bash
IMAGE = parseable/kube-collector
TAG ?= $(shell git rev-parse --short HEAD)

run: fmt vet 
	go run ./main.go

# Run go fmt against code
tidy:
	go mod tidy

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build the docker image
docker-build:
	docker build . -t ${IMAGE}:${TAG}

# Push the docker image
docker-push:
	docker push ${IMAGE}:${TAG}

# helm deploy
helm-update:
	helm upgrade --install \
	kube-collector \
	helm/kube-collector \
	-f helm/kube-collector/values.yaml --namespace kube-collector

# helm template
helm-template:
	helm template \
	kube-collector \
	helm/kube-collector \
	-f helm/kube-collector/values.yaml --namespace kube-collector