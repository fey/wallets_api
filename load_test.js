import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:8080/api/v1/wallets';

export const options = {
  vus: 1000,
  duration: '5s',
}

const getRandomOperationType = () => {
  const operations = ['DEPOSIT', 'WITHDRAW'];
  return operations[Math.floor(Math.random() * operations.length)];
}

const getRandomInt= (min, max) => {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}


export default () => {
  const walletId = getRandomInt(1, 4);
  const payload = JSON.stringify({
    WalletId: `550e8400-e29b-41d4-a716-44665544000${walletId}`,
    OperationType: getRandomOperationType(),
    Amount: getRandomInt(1, 1000),
  });
  // console.log(getRandomOperationType());
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(BASE_URL, payload, params);

  // console.log({body: res.body, status: res.status});
  check(res, {
    'is status 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
  sleep(0.001);
}
