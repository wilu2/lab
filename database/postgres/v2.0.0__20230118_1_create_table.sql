create table service (
    id serial primary key,
    name varchar,
    upstream_id varchar,
    api_set jsonb,
    creator_id integer,
    ctime timestamp without time zone not null default current_timestamp,
    last_editor_id integer,
    last_update_time timestamp without time zone not null  default current_timestamp,
    abandoned boolean default false,
    del_unique_key integer default 0
);

ALTER TABLE service ADD CONSTRAINT uniq_ser_name UNIQUE (name, del_unique_key);

create table channel (
    id serial primary key,
    name varchar,
    creator_id integer,
    ctime timestamp without time zone not null default current_timestamp,
    last_editor_id integer,
    last_update_time timestamp without time zone not null  default current_timestamp,
    abandoned boolean default false,
    del_unique_key integer default 0
);

ALTER TABLE channel ADD CONSTRAINT uniq_chn_name UNIQUE (name, del_unique_key);

create table application (
    id serial primary key,
    name varchar,
    channel_id integer,
    service_id integer,
    route_id varchar,
    status integer,
    creator_id integer,
    ctime timestamp without time zone not null default current_timestamp,
    last_editor_id integer,
    last_update_time timestamp without time zone not null  default current_timestamp,
    abandoned boolean default false,
    del_unique_key integer default 0
);

ALTER TABLE application ADD CONSTRAINT uniq_app_name UNIQUE (name, del_unique_key);

create type role as enum ('admin', 'user', 'view');

create table users (
    id serial primary key,
    account varchar,
    password varchar,
    alias varchar,
    role role,
    channels jsonb,
    description varchar,
    ctime timestamp without time zone not null default current_timestamp,
    last_update_time timestamp without time zone not null  default current_timestamp,
    abandoned boolean default false,
    del_unique_key integer default 0
);

ALTER TABLE users ADD CONSTRAINT uniq_user_account UNIQUE (account, del_unique_key);

create table config (
    id   integer,
    name varchar,
    logo varchar
);

create table session (
    id serial primary key,
    user_id integer not null,
    session varchar not null,
    expiry timestamp without time zone not null,
    abandoned boolean default false
);