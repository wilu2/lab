-- alter type role add value 'switch';

create table document (
    id serial primary key,
    content varchar
);

alter table service add document_id integer default 0;