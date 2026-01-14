-- +goose Up
-- +goose StatementBegin
-- 升级：创建表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- 降级：删除表
DROP TABLE users;
-- +goose StatementEnd

-- 00001_ is to ensure order of the migration files. Some prefer date/time prefixes like YYYYMMDDHHMM_
