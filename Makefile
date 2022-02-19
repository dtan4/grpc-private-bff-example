SERVICE := grpc-private-bff-example

PROJECT_ID ?=
REGION := asia-northeast1

TAG := $(shell git rev-parse --short HEAD)
IMAGE := gcr.io/$(PROJECT_ID)/$(SERVICE)/server:$(TAG)

.PHONY: build-server-image
build-server-image:
	docker build -t $(IMAGE) -f dockerfiles/server/Dockerfile .

.PHONY: run-server-image
run-server-image:
	docker run --rm -p 8080:8080 $(IMAGE)

.PHONY: push-server-image
push-server-image:
	docker push $(IMAGE)

.PHONY: deploy-server
deploy-server: build-server-image push-server-image
	gcloud run deploy $(SERVICE)-server \
		--project $(PROJECT_ID) \
		--region $(REGION) \
		--port 8080 \
		--image $(IMAGE) \
		--no-allow-unauthenticated \
		--use-http2
