CREATE TABLE users (
    id SERIAL PRIMARY KEY NOT NULL,
    username text NOT NULL UNIQUE,
    password text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    deleted_at timestamp
);
