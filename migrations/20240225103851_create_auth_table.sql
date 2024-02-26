-- +goose Up
create table auth (
    id serial primary key,
    name text unique not null,
    password varchar not null,
    password_confirm varchar not null,
    email text unique not null,
    role smallint not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table auth;