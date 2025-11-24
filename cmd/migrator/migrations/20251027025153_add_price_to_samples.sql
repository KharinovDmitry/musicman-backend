-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE samples ADD COLUMN price INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE samples DROP COLUMN price;
-- +goose StatementEnd

