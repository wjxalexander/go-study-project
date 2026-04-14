-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY, -- stores raw binary data (e.g. token hash), maps to []byte in Go
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL, -- TIMESTAMP(0): precision, no milliseconds
    scope TEXT NOT NULL -- "authentication" or "refresh"
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens;
-- +goose StatementEnd
