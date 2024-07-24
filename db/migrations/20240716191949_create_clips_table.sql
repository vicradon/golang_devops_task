-- +goose Up
-- +goose StatementBegin
CREATE TABLE clips (
    id SERIAL PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    content TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE clips;
-- +goose StatementEnd