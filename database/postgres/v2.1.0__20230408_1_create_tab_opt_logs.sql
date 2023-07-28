create table optlogs (
    id serial primary key,
    opt_time timestamp without time zone not null default current_timestamp,
    operation varchar,
    resource varchar,
    resource_type varchar,
    user_id integer not null,
    req_body varchar
);