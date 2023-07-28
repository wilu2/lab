create table versions (
    id serial primary key,
    service_id integer,
    version varchar,
    upstream_map jsonb default '[]'
);

ALTER TABLE versions ADD CONSTRAINT uniq_ser_version UNIQUE (service_id, version);