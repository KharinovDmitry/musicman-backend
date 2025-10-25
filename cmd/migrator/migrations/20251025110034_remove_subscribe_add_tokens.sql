-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users DROP COLUMN subscribe;
ALTER TABLE users ADD COLUMN tokens integer default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users DROP COLUMN tokens;
ALTER TABLE users ADD COLUMN subscribe boolean default false;
-- +goose StatementEnd
