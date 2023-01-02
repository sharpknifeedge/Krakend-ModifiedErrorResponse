#!make

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=pg
DOCKER_PATH=./docker/app

all: test build
build:
		$(GOBUILD) -o ./cmd/$(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) ./cmd/main.go
		./$(BINARY_NAME)
deps:
		$(GOMOD) tidy
		$(GOMOD) vendor

# SQLBoiler DB schema generator
db-schema:
		sqlboiler -c ./schema/sqlboiler.toml -o schema/entities -p entities --wipe mysql
		rm -rf ./schema/entities/*_test.go

# Docker compilation
docker-build:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DOCKER_PATH) -v
		docker-compose down
		docker-compose up --build