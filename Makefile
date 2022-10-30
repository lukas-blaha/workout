WORKOUT_BINARY=workoutApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_workout
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_workout: builds the workout binary as a linux executable
build_workout:
	@echo "Building workout binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${WORKOUT_BINARY} ./cmd/api
	@echo "Done!"
