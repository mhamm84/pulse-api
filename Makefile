# Include variables from the .envrc file
include .envrc

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## confirm: create the new confirm target.
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## api/audit: tidy dependencies and format, vet and test all code
.PHONY: api/audit
api/audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...

## api/build: build local go binary and linux_amd_64 binary
.PHONY: api/build
api/build: api/audit
	@echo "Building pulse API..."
	go build ./cmd/pulse/
	GOOS=linux GOARCH=amd64 go build -o=./docker/bin/linux_amd64/pulse ./cmd/pulse/

## api/go/run: run the cmd/api application
.PHONY: api/go/run
api/go/run:
	go run ./cmd/pulse/ run-api --db-dsn=${PULSE_POSTGRES_DSN} --cors-trusted-origins="http://localhost:9090" --log-level="DEBUG"

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new: confirm
	@echo 'Create a migration for file ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path=./migrations -database=${PULSE_POSTGRES_FROM_HOST_DSN} up

## db/migrations/down: apply all up database migrations
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Running down migrations...'
	migrate -path=./migrations -database=${PULSE_POSTGRES_FROM_HOST_DSN} down

## db/migrations/force version=$1: on any migration error, force to the version passed in
.PHONY: db/migrations/force
db/migrations/force:
	@echo 'Running migration force for version: ${version}...'
	migrate -path=./migrations -database=${PULSE_POSTGRES_FROM_HOST_DSN} force ${version}

## dev/docker/build: build docker image for the pulse api
.PHONY: dev/docker/build
dev/docker/build: api/build
	@echo "Building dev Docker env"
	go mod vendor
	docker-compose -f docker-compose-dev.yml build

## dev/docker/up: create and start Pulse API containers for dev
.PHONY: dev/docker/up
dev/docker/up: dev/docker/build
	@echo "Creating and starting Pulse API containers..."
	docker-compose -f docker-compose-dev.yml up -d

## dev/docker/start: start all docker containers for dev
.PHONY: dev/docker/start
dev/docker/start:
	@echo "Starting Pulse API stack..."
	docker-compose -f docker-compose-dev.yml start

## dev/docker/stop: stop all docker containers
.PHONY: dev/docker/stop
dev/docker/stop:
	@echo "Stopping Pulse API stack..."
	docker-compose -f docker-compose-dev.yml stop

## dev/pulse/down: stop and remove all docker containers
.PHONY: dev/docker/down
dev/docker/down:
	@echo "Stopping & removing Pulse API containers..."
	docker-compose -f docker-compose-dev.yml down

## integration/up: create & start integration test docker containers
.PHONY: integration/up
integration/up: dev/docker/build
	@echo 'Spinning up docker containers for integration tests...'
	cd docker/integration ; \
		docker-compose up -d

## integration/tests/run: run integration tests
.PHONY: integration/tests/run
integration/tests/run: api/audit integration/up
	@echo 'Running tests...'
	go test -v -tags=integration -race ./...

## integration/down: spin down the integration docker stack
.PHONY: integration/down
integration/down:
	@echo 'Spinning up docker for integration tests'
	cd ./docker/integration ; \
		docker-compose down

## unit/tests/run: run all unit tests
.PHONY: unit/tests/run
unit/tests/run: api/audit
	@echo 'Running unit tests...'
	go test -v -race ./...