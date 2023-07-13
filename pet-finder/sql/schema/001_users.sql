-- +goose Up
CREATE TYPE roles as ENUM ('admin','user');
CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY,
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL,
    name text NOT NULL,
    email text NOT NULL UNIQUE,
    password text NOT NULL,
    role roles NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT encode(sha256(random()::text::bytea), 'hex')



);

-- +goose Down
DROP TABLE users;
DROP TYPE roles;
