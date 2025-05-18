CREATE TABLE IF NOT EXISTS balances (
    user_id UUID PRIMARY KEY,
    available DECIMAL NOT NULL CHECK(available >= 0),
    withdrawn DECIMAL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
