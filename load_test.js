import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:8080/api/v1/wallets';

export let options = {
  vus: 1000,
  duration: '1s',
}

export default function () {
  const payload = JSON.stringify({
    WalletId: '550e8400-e29b-41d4-a716-446655440000', // Замените на актуальный ID кошелька
    OperationType: 'DEPOSIT', // Или 'Withdraw'
    Amount: 100, // Сумма операции
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  let res = http.post(BASE_URL, payload, params);

  // Проверка статуса ответа
  check(res, {
    'is status 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });

  sleep(0.001);
}
