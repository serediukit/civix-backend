CREATE TABLE users
(
    user_id       serial primary key,
    email         varchar(255) not null,
    password_hash varchar(255) not null,
    name          varchar(255) not null,
    created_at    timestamp    not null,
    updated_at    timestamp    not null,
    deleted_at    timestamp
);

CREATE UNIQUE INDEX users_email_unique_index ON users (email);