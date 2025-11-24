-- +goose Up
-- +goose StatementBegin
CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    sample_id UUID NOT NULL REFERENCES samples(id) ON DELETE CASCADE,
    price INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_uuid, sample_id)
);

CREATE INDEX idx_purchases_user_uuid ON purchases(user_uuid);
CREATE INDEX idx_purchases_sample_id ON purchases(sample_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_purchases_sample_id;
DROP INDEX IF EXISTS idx_purchases_user_uuid;
DROP TABLE IF EXISTS purchases;
-- +goose StatementEnd

