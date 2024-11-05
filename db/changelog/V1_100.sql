create table if not exists wombatsm.message_event
(
    hash        uuid primary key,
    source_type varchar(64)              not null,
    text        text                     not null,
    author_id   varchar                  not null,
    chat_id     varchar                  not null,
    message_id  varchar                  not null,
    comment_id  varchar                  not null,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp
);

create table if not exists wombatsm.acl
(
    author_id  varchar primary key,
    is_allowed boolean not null
)
