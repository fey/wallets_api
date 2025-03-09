install:
	go get .

start:
	go run main.go

test:
	go test

lint:
	# TODO: add golangci-lint

build:
	go build -o wallets_api .

compose-start:
	docker compose up --abort-on-container-failure

compose-stop:
	docker compose down

compose-build:
	docker compose build

compose-bash:
	docker compose run --rm app bash

compose-logs:
	docker compose logs -f

compose-test:
	docker compose -f docker-compose.test.yml -f docker-compose.yml -p wallets_api-tests up

compose-production-build:
	docker compose -f docker-compose.production.yml build

compose-production-start:
	docker compose -f docker-compose.production.yml up --abort-on-container-failure --build
