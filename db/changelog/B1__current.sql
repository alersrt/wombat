create extension if not exists "uuid-ossp";

create schema if not exists wombatsm;

create table if not exists wombatsm.accounts
(
    gid       uuid primary key                  default gen_random_uuid(),
    create_ts timestamp with time zone not null default current_timestamp,
    update_ts timestamp with time zone not null default current_timestamp
);

create table if not exists wombatsm.source_connections
(
    gid         uuid primary key                  default gen_random_uuid(),
    account_gid uuid                     not null,
    source_type varchar(64)              not null,
    user_id     varchar                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp,
    constraint source_connections_account_gid_fk foreign key (account_gid) references wombatsm.accounts (gid),
    constraint source_connections_nk unique (source_type, user_id)
);

create table if not exists wombatsm.target_connections
(
    gid         uuid primary key                  default gen_random_uuid(),
    account_gid uuid                     not null,
    target_type varchar(64)              not null,
    token       bytea                    not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp,
    constraint target_connections_account_gid_fk foreign key (account_gid) references wombatsm.accounts (gid),
    constraint target_connections_nk unique (target_type, account_gid)
);

create table if not exists wombatsm.comments
(
    gid         uuid primary key                  default gen_random_uuid(),
    target_type varchar(64)              not null,
    source_type varchar(64)              not null,
    comment_id  varchar                  not null,
    user_id     varchar                  not null,
    chat_id     varchar                  not null,
    message_id  varchar                  not null,
    tag         varchar                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp,
    constraint comments_nk unique (target_type, comment_id)
);
