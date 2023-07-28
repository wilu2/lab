alter table users add group_id integer not null DEFAULT -1;

ALTER table users add sso_code integer not null DEFAULT 0;

ALTER table application add group_id integer NOT NULL DEFAULT 0;

UPDATE users set sso_code = 1 where id = 1;

ALTER table users drop constraint uniq_user_account;

ALTER table users add constraint uniq_user_account unique (sso_code, account, del_unique_key);

Alter table service add group_id integer not null DEFAULT -1;

Alter table channel add group_id integer not null DEFAULT -1;

ALTER table channel drop constraint uniq_chn_name;

ALTER table channel add constraint uniq_chn_name unique (name, group_id, del_unique_key);

create type sv_type as enum ('extract', 'sort', 'ext&sort', 'others');
-- extract - 提取 
-- sort - 分类 
-- ext&sort - 抽取+分类 
-- others - 其它

Alter table service add service_type sv_type not null DEFAULT 'others';


