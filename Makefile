PROJECT_ID ?= turnero-165d4
REGION ?= us-central1
SERVICE_NAME ?= ServiceBookingApp
IMAGE_NAME ?= gcr.io/$(PROJECT_ID)/$(SERVICE_NAME)

.PHONY: run build test docker-build docker-push deploy

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test ./...

docker-build:
	docker build -t $(IMAGE_NAME) .

docker-push:
	docker push $(IMAGE_NAME)

deploy:
	gcloud run deploy $(SERVICE_NAME) \
		--image $(IMAGE_NAME) \
		--region $(REGION) \
		--platform managed \
		--allow-unauthenticated \
		--project $(PROJECT_ID)
