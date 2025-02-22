create extension if not exists "uuid-ossp";
create schema if not exists wombatsm;

create table if not exists wombatsm.comments
(
    gid         uuid primary key                  default gen_random_uuid(),
    comment_id  varchar                  not null,
    target_type varchar(64)              not null,
    text        text                     not null,
    source_type varchar(64)              not null,
    author_id   varchar                  not null,
    chat_id     varchar                  not null,
    message_id  varchar                  not null,
    tag         varchar                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp,
    constraint comments_comment_target_u unique (comment_id, target_type)
);

create table if not exists wombatsm.acl
(
    gid         uuid primary key                  default gen_random_uuid(),
    author_id   varchar                  not null,
    source_type varchar(64)              not null,
    is_allowed  boolean                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp,
    constraint acl_author_source_u unique (author_id, source_type)
);

create table if not exists wombatsm.connections
(
    gid         uuid primary key                  default gen_random_uuid(),
    author_id   varchar                  not null,
    source_type varchar(64)              not null,
    target_type varchar(64)              not null,
    token       varchar                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp,
    constraint acl_author_source_fk foreign key (author_id, source_type) references wombatsm.acl (author_id, source_type)
)
