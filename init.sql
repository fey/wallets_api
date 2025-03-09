DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;

CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    wallet_id UUID REFERENCES wallets(id),
    operation_type VARCHAR(8),
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO wallets (id, balance, created_at, updated_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 100.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440001', 250.50, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440002', 75.25, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440003', 500.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440004', 0.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO transactions (wallet_id, operation_type, amount, created_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 'deposit', 50.00, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440001', 'withdraw', 20.00, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440002', 'deposit', 100.00, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440003', 'withdraw', 10.00, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440004', 'deposit', 200.00, CURRENT_TIMESTAMP);
