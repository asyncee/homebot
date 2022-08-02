BINARY_NAME=homebot
BUILD_DIR=build

.PHONY: pre-build
pre-build:
	@mkdir -p $(BUILD_DIR)

.PHONY: build-linux
build-linux: pre-build
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 cmd/homebot/main.go

.PHONY: build-osx
build-osx: pre-build
	@echo "Building OSX binary..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 cmd/homebot/main.go

.PHONY: build
build: clean build-linux build-osx

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -Rf $(BUILD_DIR)

.PHONY: run
run:
	go run cmd/homebot/main.go

.PHONY: docker-build
docker-build:
	docker build --tag homebot:$$(git describe --tags) .
	docker tag homebot:$$(git describe --tags) asyncee/homebot:$$(git describe --tags)
	docker tag homebot:$$(git describe --tags) asyncee/homebot:latest
	docker push asyncee/homebot -a

.PHONY: docker-run
docker-run:
	docker run --rm --env-file=.env asyncee/homebot
