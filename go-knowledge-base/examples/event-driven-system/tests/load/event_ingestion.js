import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 500 },
    { duration: '5m', target: 1000 },
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'],
    http_req_failed: ['rate<0.01'],
  },
};

function createOrder() {
  const customerId = uuidv4();
  const orderId = uuidv4();
  
  const payload = JSON.stringify({
    customer_id: customerId,
    items: [
      {
        product_id: uuidv4(),
        product_name: 'Test Product 1',
        quantity: Math.floor(Math.random() * 10) + 1,
        unit_price: 29.99
      },
      {
        product_id: uuidv4(),
        product_name: 'Test Product 2',
        quantity: Math.floor(Math.random() * 5) + 1,
        unit_price: 49.99
      }
    ]
  });

  const res = http.post('http://localhost:8080/api/v1/orders', payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'order created': (r) => r.status === 201,
    'has order id': (r) => r.json('order_id') !== undefined,
  });

  return res.json('order_id');
}

function payOrder(orderId) {
  const payload = JSON.stringify({
    payment_id: uuidv4(),
    amount: 100.00,
    payment_method: 'credit_card'
  });

  const res = http.post(`http://localhost:8080/api/v1/orders/${orderId}/pay`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'payment processed': (r) => r.status === 200,
  });
}

function queryOrder(orderId) {
  const res = http.get(`http://localhost:8080/api/v1/orders/${orderId}`);
  
  check(res, {
    'order retrieved': (r) => r.status === 200,
    'has correct structure': (r) => r.json('order_id') === orderId,
  });
}

export default function () {
  // Create order
  const orderId = createOrder();
  
  sleep(0.5);
  
  // Query order
  queryOrder(orderId);
  
  sleep(0.5);
  
  // Pay order (50% of orders)
  if (Math.random() > 0.5) {
    payOrder(orderId);
  }
  
  sleep(1);
}
