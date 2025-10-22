CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_expiry ON tokens(expiry);

-- Token scopes: 'authentication', 'activation'