.PHONY: help test test-race test-leak bench bench-compare lint sec-scan 

help: ## show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

PROJECT_NAME?=kol-resource
SHELL = /bin/bash

########
# test #
########

test: test-race test-leak ## launch all tests

test-race: ## launch all tests with race detection
	go test ./... -cover -race

test-leak: ## launch all tests with leak detection (if possible)
	go test ./... -leak

test-coverage-report:
	go test -v  ./... -cover -race -covermode=atomic -coverprofile=./coverage.out
	go tool cover -html=coverage.out


#############
# benchmark #
#############

bench: ## launch benchs
	go test ./... -bench=. -benchmem | tee ./bench.txt

bench-compare: ## compare benchs results
	benchstat ./bench.txt

########
# lint #
########

lint: ## lints the entire codebase
	@golangci-lint run ./... --config=./.golangci.toml

#######
# sec #
#######

sec-scan: trivy-scan vuln-scan ## scan for security and vulnerability issues

trivy-scan: ## scan for sec issues with trivy (trivy binary needed)
	trivy fs --exit-code 1 --no-progress --severity CRITICAL ./

vuln-scan: ## scan for vulnerability issues with govulncheck (govulncheck binary needed)
	govulncheck ./...


#######
# sql #
#######

sqlboiler:
	@( \
	printf "Enter pass for db: "; read -s DB_PASSWORD && \
	printf "Enter port(5432, 26257...): \n"; read -r DB_PORT && \
	PSQL_HOST=localhost PSQL_PORT=$$DB_PORT PSQL_PASS=$$DB_PASSWORD PSQL_DBNAME=$(APP_NAME_UND) PSQL_USER=root sqlboiler psql -c ./database/sqlboiler/sqlboiler.toml && \
	go get -t kolresource/internal/db/sqlboiler && \
	PSQL_HOST=localhost PSQL_PORT=$$DB_PORT PSQL_PASS=$$DB_PASSWORD PSQL_DBNAME=$(APP_NAME_UND) PSQL_USER=root sqlboiler psql -c ./database/sqlboiler/sqlboiler.toml --no-tests=true && \
	go mod tidy \
	)

migration-up:
	@( \
	printf "Enter pass for db: \n"; read -s DB_PASSWORD && \
	printf "Enter port(5432, 26257...): \n"; read -r DB_PORT &&\
	migrate -database "postgres://root:$${DB_PASSWORD}@localhost:$${DB_PORT}/$(APP_NAME_UND)?sslmode=disable" -path database/migrations up \
	)

migration-down:
	@( \
	printf "Enter pass for db: \n"; read -s DB_PASSWORD && \
	printf "Enter port(5432, 26257...): \n"; read -r DB_PORT &&\
	migrate -database "postgres://root:$${DB_PASSWORD}@127.0.0.1:$${DB_PORT}/$(APP_NAME_UND)?sslmode=disable" -path database/migrations down \
	)

SQL_FILE_TIMESTAMP=$(shell date '+%Y%m%d%H%M%S')

gen-migrate-sql:
	@( \
	printf "Enter file name: "; read -r FILE_NAME; \
	touch database/migrations/$(SQL_FILE_TIMESTAMP)_$$FILE_NAME.up.sql; \
	touch database/migrations/$(SQL_FILE_TIMESTAMP)_$$FILE_NAME.down.sql; \
	)


