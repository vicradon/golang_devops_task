-- +goose Up
-- +goose StatementBegin
ALTER TABLE clips ADD COLUMN delete_after TIMESTAMP NULL;

-- Set delete_after to a far future date where it's intended to be viewed once and deleted
UPDATE clips
SET
    delete_after = '9999-12-31 23:59:59'
WHERE
    delete_after IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE clips DROP COLUMN delete_after;
-- +goose StatementEnd