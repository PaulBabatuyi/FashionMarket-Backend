-- +goose Down
ALTER TABLE users
  DROP COLUMN IF EXISTS address,
  DROP COLUMN IF EXISTS country;