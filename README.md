

```bash
curl -X POST http://localhost:8080/api/v1/wallet -H "Content-Type: application/json" -d '{
    "walletId": "123e4567-e89b-12d3-a456-426614174000",
    "operationType": "DEPOSIT",
    "amount": 1000
}'
```

```bash
curl -X POST http://localhost:8080/api/v1/wallet -H "Content-Type: application/json" -d '{
    "walletId": "123e4567-e89b-12d3-a456-426614174000",
    "operationType": "WITHDRAW",
    "amount": 500
}'
```

```bash
{
    "walletId": "123e4567-e89b-12d3-a456-426614174000",
    "balance": 500
}
```
