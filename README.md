# wallets_api

REST API приложение, которое работает с базой данных и позволяет управлять балансом электронного кошелька.

## Requirements

* Ubuntu (Linux)
* Make
* Go 1.12.1+
* Docker
* k6 for testing

## Usage

```bash
docker pull ghcr.io/fey/wallets_api:latest
```

Start:

```bash
docker run \
    --rm \
    --name wallets_api \
    -p 8080:8080 \
    -v ($pwd)/path/to/config:/app/config.env
    ghcr.io/fey/wallets_api
```

For Docker Compose see [Makefile](./Makefile)

## API

## Wallet

```bash
curl localhost:8080/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000
```

Response:

```json
{"wallet_id":"550e8400-e29b-41d4-a716-446655440000","balance":1100}
```

### Operations

Deposit:

```bash
curl -X POST http://localhost:8080/api/v1/wallets -H "Content-Type: application/json" -d '{
    "walletId": "550e8400-e29b-41d4-a716-446655440000",
    "operationType": "DEPOSIT",
    "amount": 1000
}'
```

Withdraw:

```bash
curl -X POST http://localhost:8080/api/v1/wallets -H "Content-Type: application/json" -d '{
    "walletId": "550e8400-e29b-41d4-a716-446655440000",
    "operationType": "WITHDRAW",
    "amount": 1000
}'
```

Response:

```json
{"wallet_id":"550e8400-e29b-41d4-a716-446655440000","balance":2100}
```

For more see [swagger.yaml](./docs/swagger.yaml) or `<app_url>/swagger/index.html`
