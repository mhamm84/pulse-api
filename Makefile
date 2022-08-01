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

## run/api: run the cmd/api application
.PHONY: api/run
api/run:
	go run ./cmd/api/ -db-dsn=${PULSE_POSTGRES_DSN} -cors-trusted-origins="http://localhost:9090" -log-level="INFO"

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new: confirm
	@echo 'Create a migration for file ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running migrations...'
	migrate -path=. -database=${PULSE_POSTGRES_DSN} up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...

.PHONY: api/build
api/build:
	@echo "Building pulse API..."
	go build -o=./docker/bin/api ./cmd/api/
	GOOS=linux GOARCH=amd64 go build -o=./docker/bin/linux_amd64/api ./cmd/api/


.PHONY: api/docker/image
api/docker/image: api/build
	@echo "Building pulse API Docker Image"
	docker build -t mhamm84/pulse-api ./docker

.PHONY: integration-docker-up
integration-docker-up:
	@echo 'Spinning up docker for integration tests'
	cd docker/integration ; \
		docker-compose up -d

.PHONY: integration-tests
integration-tests: audit integration-docker-up
	@echo 'Running tests...'
	go test -v -tags=integration -race ./...

.PHONY: integration-docker-down
integration-docker-down:
	@echo 'Spinning up docker for integration tests'
	cd ./docker/integration ; \
		docker-compose down

.PHONY: unit-tests
unit-tests: audit
	@echo 'Running unit tests...'
	go test -v -race ./...