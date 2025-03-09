# wallets_api

REST API приложение, которое работает с базой данных и позволяет управлять балансом электронного кошелька.

## Requirements

* Ubuntu (Linux)
* Make
* Go 1.12.1+
* Docker

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

_config.env_ file example:

```dotenv
DATABASE_URL=postgres://wallets_api:password@db:5432/wallets_api_dev?sslmode=disable
```

For Docker Compose see [Makefile](./Makefile)

### Operations

Deposit:

```bash
curl -X POST http://localhost:8080/api/v1/wallet -H "Content-Type: application/json" -d '{
    "walletId": "123e4567-e89b-12d3-a456-426614174000",
    "operationType": "DEPOSIT",
    "amount": 1000
}'
```



Withdraw:

```bash
curl -X POST http://localhost:8080/api/v1/wallet -H "Content-Type: application/json" -d '{
    "walletId": "123e4567-e89b-12d3-a456-426614174000",
    "operationType": "WITHDRAW",
    "amount": 500
}'
```
