CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER NOT NULL,
    user_id UUID NOT NULL,
    order_number VARCHAR(16) NOT NULL,
    operation VARCHAR(16) NOT NULL,
    amount DECIMAL NOT NULL,
    proccessedAt TIMESTAMPTZ NOT NULL,

    PRIMARY KEY(id, order_id, proccessedAt)
);

SELECT create_hypertable(
    'transactions',
    'proccessedAt',
    chunk_time_interval => INTERVAL '6 days',
    if_not_exists => TRUE
);

CREATE TABLE IF NOT EXISTS transaction_counter (
    user_id UUID PRIMARY KEY,
    counter INTEGER NOT NULL
);
