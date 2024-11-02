create table if not exists wombatsm.message_event
(
    hash        uuid primary key         not null,
    source_type varchar(64),
    event_type  varchar(64),
    text        text,
    author_id   varchar,
    chat_id     varchar,
    message_id  varchar,
    create_ts   timestamp with time zone not null default current_timestamp,
    update_ts   timestamp with time zone not null default current_timestamp
)
