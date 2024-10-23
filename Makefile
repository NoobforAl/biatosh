#======================================================================
# Makefile for biatosh
#
# Usage:
#   make [target]
#
# Targets:
#   all          - Build and run the project
#   build        - Build the Docker image
#   run          - Run the code outside of Docker
#	dev          - Run the code outside of Docker with air
#   dev-add-create-user - Run the code outside of Docker with air and add create-user
#   run-docker   - Run the code inside Docker
#   clean        - Clean up Docker containers and images
#   help         - Display this help message
#======================================================================

DOCKER_IMAGE_NAME=biatosh:latest
DOCKER_CONTAINER_NAME=biatosh
SOURCE_DIR=src
MAIN_FILE=main.go

# Targets
.PHONY: all build run run-docker clean help

all: build run

build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) .

run: _generate_sqlc
	@echo "Running code outside of Docker..."
	cd $(SOURCE_DIR) && go run $(MAIN_FILE)

dev:
	@echo "Running code outside of Docker..."
	cd $(SOURCE_DIR) && air

dev-add-create-user: _generate_sqlc
	@echo "Enter name: " && read name && \
	echo "Enter email: " && read email && \
	echo "Enter password: " && read password && \
	echo "Enter username: " && read username && \
	echo "Enter phone: " && read phone && \
	echo "Running code outside of Docker..." && \
	cd $(SOURCE_DIR) && go run $(MAIN_FILE) create-user --name="$$name" --email="$$email" --password="$$password" --username="$$username" --phone="$$phone"

_generate_sqlc:
	@echo "Generating SQLC..."
	cd $(SOURCE_DIR) && sqlc generate

run-docker:
	@echo "Running code inside Docker..."
	docker run --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME)

clean:
	@echo "Cleaning up..."
	rm -rf $(SOURCE_DIR)/tmp
	docker rm -f $(DOCKER_CONTAINER_NAME) || true
	docker rmi -f $(DOCKER_IMAGE_NAME) || true

help:
	@echo "======================================================================"
	@echo "Makefile for MyProject"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all          - Build and run the project"
	@echo "  build        - Build the Docker image"
	@echo "  run          - Run the code outside of Docker"
	@echo "	 dev          - Run the code outside of Docker with air"
	@echo "	 dev-add-create-user - Run the code outside of Docker with air and add create-user"
	@echo "  run-docker   - Run the code inside Docker"
	@echo "  clean        - Clean up Docker containers and images"
	@echo "  help         - Display this help message"
	@echo "======================================================================"