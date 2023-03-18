drop table if exists user;

drop table if exists user_invite_record;

create table user
(
    id           INTEGER not null
        primary key autoincrement,
    user_id      TEXT    not null,
    user_name    text,
    remain_count INTEGER default 0,
    invite_code  TEXT
);

create index idx_invite_link
    on user (invite_code);

create unique index idx_user_id
    on user (user_id);

create table user_invite_record
(
    id             integer not null
        primary key autoincrement,
    user_id        text    not null,
    invite_user_id text    not null,
    invite_time    text
);

create index idx_invite_record_user_id
    on user_invite_record (user_id);

create index idx_invite_record_invite_time
    on user_invite_record (invite_time);

