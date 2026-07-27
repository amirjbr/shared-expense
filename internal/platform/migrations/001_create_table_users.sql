-- +migrate Up
CREATE TABLE users (
    id UUID PRIMARY KEY ,
    first_name VARCHAR,
    last_name VARCHAR,
    username VARCHAR,
    password VARCHAR,
    email VARCHAR,
    phone_number VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +migrate Down
DROP TABLE users;