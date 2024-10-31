create table if not exists wombatsm.message_event
(
    hash        uuid,
    source_type varchar(64),
    event_type  varchar(64),
    text        text,
    author_id   varchar,
    chat_id     varchar,
    message_id  varchar
)
