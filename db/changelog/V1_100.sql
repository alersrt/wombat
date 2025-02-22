create extension if not exists "uuid-ossp";
create schema if not exists wombatsm;

create table if not exists wombatsm.comments
(
    comment_id  varchar primary key      not null,
    target_type varchar(64)              not null,
    source_type varchar(64)              not null,
    text        text                     not null,
    author_id   varchar                  not null,
    chat_id     varchar                  not null,
    message_id  varchar                  not null,
    tag         varchar                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp
);

create table if not exists wombatsm.acl
(
    author_id   varchar                  not null,
    source_type varchar(64)              not null,
    is_allowed  boolean                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp,
    primary key (author_id, source_type)
);

create table if not exists wombatsm.connections
(
    author_id   varchar references acl (author_id),
    source_type varchar(46) references acl (source_type),
    target_type varchar(64)              not null,
    token       varchar                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp
)
