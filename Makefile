DOCKER_IMAGE_NAME = news-aggregator
DOCKER_IMAGE_VERSION = 1.0

# Build the cli application
build:
	go build -o bin/news-aggregator cmd/cli/main/main.go

# Build the docker image and run the container
build-and-run-docker-image:
	$(MAKE) build-docker-image
	$(MAKE) run-docker-image

# Stop the docker container and remove it
stop-and-remove-docker-container:
	$(MAKE) stop-docker-container
	$(MAKE) remove-container

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