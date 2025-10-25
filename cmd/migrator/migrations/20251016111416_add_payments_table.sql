-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists payments(
   id  VARCHAR(64) NOT NULL,
   user_uuid UUID NOT NULL,
   payment_status VARCHAR(50) NOT NULL,
   description TEXT,
   amount INT NOT NULL,
   captured_at TIMESTAMP,
   created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists payments;
-- +goose StatementEnd
