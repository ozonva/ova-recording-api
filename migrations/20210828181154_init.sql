-- +goose Up
-- +goose StatementBegin
CREATE TABLE appointments (
  appointment_id SERIAL PRIMARY KEY,
  user_id bigint NOT NULL,
  name text NOT NULL,
  description text NOT NULL DEFAULT '',
  start_time date NOT NULL,
  end_time date NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE appointments;
-- +goose StatementEnd
