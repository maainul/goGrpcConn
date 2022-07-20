-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS blogs
(

    id               VARCHAR(100) PRIMARY KEY DEFAULT uuid_generate_v4(),
    name             VARCHAR(255) NOT NULL DEFAULT '',
    created_at       TIMESTAMP    DEFAULT current_timestamp,
    created_by       VARCHAR(100) NOT NULL DEFAULT '',
    updated_at       TIMESTAMP    DEFAULT current_timestamp,
    updated_by       VARCHAR(100) NOT NULL DEFAULT '',
    deleted_at       TIMESTAMP    DEFAULT NULL,
    deleted_by       VARCHAR(100) NOT NULL DEFAULT ''
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS blogs;
