SERVICE := grpc-private-bff-example

PROJECT_ID ?=
REGION := asia-northeast1

TAG := $(shell git rev-parse --short HEAD)
IMAGE := gcr.io/$(PROJECT_ID)/$(SERVICE)/server:$(TAG)

CLOUD_RUN_ENDPOINT ?=

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

.PHONY: build-bff-image
build-bff-image:
	docker build -t bff -f dockerfiles/bff/Dockerfile .

.PHONY: run-bff-image
run-bff-image:
	docker run \
		--rm \
		-p 8080:8080 \
		-e CLOUD_RUN_ENDPOINT=$(CLOUD_RUN_ENDPOINT) \
		-e IMPERSONATE_SA_EMAIL=$(IMPERSONATE_SA_EMAIL) \
		-v "$${HOME}/.config/gcloud:/root/.config/gcloud:ro" \
		bff
