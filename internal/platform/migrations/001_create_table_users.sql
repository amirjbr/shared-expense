-- +migrate Up
CREATE TABLE users (
    id UUID PRIMARY KEY ,
    first_name VARCHAR,
    last_name VARCHAR,
    username VARCHAR UNIQUE ,
    password VARCHAR,
    email VARCHAR UNIQUE ,
    phone_number VARCHAR UNIQUE ,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +migrate Down
DROP TABLE users;