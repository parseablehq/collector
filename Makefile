SHELL := /bin/bash
IMAGE = parseable/collector
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
	collector \
	helm/collector \
	-f helm/collector/values.yaml \
	--create-namespace \
	--namespace collector

# helm template
helm-template:
	helm template \
	collector \
	helm/collector \
	-f helm/collector/values.yaml --namespace collector
