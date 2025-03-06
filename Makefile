start:
	go run main.go

build:
	go build -o wallets_api

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
