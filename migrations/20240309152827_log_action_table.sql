-- +goose Up
create table log_action (
    id serial primary key,
    action text not null,
    user_id integer not null
);
-- +goose Down
drop table log_action;