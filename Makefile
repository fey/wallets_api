install:
	go get .

start:
	go run main.go

test:
	go test -v

lint:
	# TODO: add golangci-lint

prepare-config:
	cp config.env.example config.env

compose-setup: prepare-config compose-build

test-load:
	# TODO
	# ./k6 load_test.js

build:
	go build -buildvcs=false -o wallets_api .

compose-start:
	docker compose up --abort-on-container-failure

compose-stop:
	docker compose down

compose-down:
	docker compose down -v --remove-orphans

compose-build:
	docker compose build

compose-bash:
	docker compose run --rm app bash

compose-logs:
	docker compose logs -f

compose-test:
	docker compose run --rm app make test

compose-production-build:
	docker compose -f docker-compose.production.yml build

compose-production-start:
	docker compose -f docker-compose.production.yml up --abort-on-container-failure --build
