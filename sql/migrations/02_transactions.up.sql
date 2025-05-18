CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER NOT NULL,
    user_id UUID NOT NULL,
    order_number VARCHAR(16),
    operation VARCHAR(16) NOT NULL,
    amount DECIMAL NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,

    PRIMARY KEY(id, order_id, timestamp)
);

SELECT create_hypertable(
    'transactions',
    'timestamp',
    chunk_time_interval => INTERVAL '6 days',
    if_not_exists => TRUE
);

CREATE TABLE IF NOT EXISTS transaction_counter (
    user_id UUID PRIMARY KEY,
    counter INTEGER NOT NULL
);
