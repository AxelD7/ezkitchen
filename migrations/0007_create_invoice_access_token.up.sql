CREATE TABLE invoice_access_tokens (
    invoice_token_id BIGSERIAL PRIMARY KEY,

    estimate_id INTEGER NOT NULL
        REFERENCES estimates(estimate_id)
        ON DELETE CASCADE,

    token_hash CHAR(64) NOT NULL,

    expires_at TIMESTAMPTZ NOT NULL,

    used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX invoice_access_tokens_token_hash_idx
    ON invoice_access_tokens (token_hash);

CREATE INDEX invoice_access_tokens_estimate_id_idx
    ON invoice_access_tokens (estimate_id);
