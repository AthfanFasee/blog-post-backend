include app.env
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]






# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_DSN} up





# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and verify dependencies
.PHONY: verify
verify:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify





# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell Get-Date -Format "MM/dd/yyyy:HH:mm")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'


## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

## dockercompose: run dockercompose up
.PHONY: dockerup
dockerup:
	@echo 'Running dockercompose up'
	docker-compose build --build-arg VERSION=${git_description} --build-arg CURRENT_TIME=${current_time}
	docker compose up

## dockercompose: run dockercompose down
.PHONY: dockerdown
dockerdown:
	@echo 'Running dockercompose down'
	docker compose down




# ==================================================================================== #
# TEST
# ==================================================================================== #

.PHONY: test/api
test/api:
	@echo 'testing cmd/api...'
	go test -v ./cmd/api/

.PHONY: test/api/race
test/api/race:
	@echo 'testing with race detector cmd/api...'
	go test -v -race ./cmd/api/




# ==================================================================================== #
# Windows
# ==================================================================================== #
.PHONY: db/up/win
db/up/win:
	@echo 'running db/migrate/up...'
	migrate -path ./migrations -database ${DB_DSN} -verbose up

.PHONY: db/down/win
db/down/win:
	@echo 'running db/migrate/down...'
	migrate -path ./migrations -database ${DB_DSN} down

.PHONY: db/downversion/win
db/downversion/win:
	@echo 'running db/migrate/force...'
	migrate -path ./migrations -database ${DB_DSN} force 1

