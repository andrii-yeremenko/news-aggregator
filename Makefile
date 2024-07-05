DOCKER_IMAGE_NAME = news-aggregator
DOCKER_IMAGE_VERSION = 1.0

# Run the project
run-web-server:
	go run cmd/web-server/main/main.go

# Build the docker image
build-docker-image:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION) .

# Run the docker image in container
run-docker-image:
	docker run -p 8080:8443 $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION)

# Stop the docker container
stop-docker-container:
	docker stop $(shell docker ps -a -q --filter ancestor=$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION) --format="{{.ID}}")

# Remove the docker container
remove-container:
	docker rm $(shell docker ps -a -q --filter ancestor=$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION) --format="{{.ID}}")