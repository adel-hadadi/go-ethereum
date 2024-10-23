APP_BINARY=webApp

build:
	@echo "Building app binary"
	env GOOS=linux CGO_ENABLED=0 go build -o ${APP_BINARY} ./cmd/api
	@echo "App binary buited"

up-build: build
	@echo "Stopping docker images"
	docker compose down
	@echo "Building and starting docker images"
	docker compose up -d --build
	@echo "Docker images built and started!"

