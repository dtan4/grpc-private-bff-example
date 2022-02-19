PROJECT_ID ?=

TAG := $(shell git rev-parse --short HEAD)
IMAGE := gcr.io/$(PROJECT_ID)/grpc-private-bff-example/server:$(TAG)

.PHONY: build-server-image
build-server-image:
	docker build -t $(IMAGE) -f dockerfiles/server/Dockerfile .

.PHONY: run-server-image
run-server-image:
	docker run --rm -p 8080:8080 $(IMAGE)
