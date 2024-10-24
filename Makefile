APP_BINARY=webApp

build:
	@echo "Building app binary"
	env GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o ${APP_BINARY} ./internal/ethereum/
	@echo "App binary built"

up-build: build
	@echo "Stopping docker images"
	docker compose down
	@echo "Building and starting docker images"
	docker compose up -d --build
	@echo "Docker images built and started!"

down:
	@echo "Stopping docker images"
	docker compose down
	@echo "Done!"

logs:
	docker compose logs -f
