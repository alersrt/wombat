package dao

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
	"wombat/internal/domain"
)

type MessageEventRepository struct {
	PostgreSQLQueryHelper[domain.MessageEvent, string]
}

func NewMessageEventRepository(url *string) (*MessageEventRepository, error) {
	db, err := sqlx.Connect("postgres", *url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &MessageEventRepository{
		PostgreSQLQueryHelper[domain.MessageEvent, string]{
			db:            db,
			entityFactory: &MessageEventEntityFactory{},
		},
	}, nil
}

func (receiver *MessageEventRepository) GetById(id string) (*domain.MessageEvent, error) {
	entity, err := receiver.GetEntityById("select * from wombatsm.message_event where hash = $1", id)
	if err != nil {
		return nil, err
	}

	return entity.ToDomain(), nil
}

func (receiver *MessageEventRepository) Save(domain *domain.MessageEvent) (*domain.MessageEvent, error) {
	query := `insert into wombatsm.message_event(hash, source_type, event_type, text, author_id, chat_id, message_id)
               values (:hash, :source_type, :event_type, :text, :author_id, :chat_id, :message_id)
               on conflict (hash)
               do update
               set event_type = :event_type,
                   text = :text,
                   author_id = :author_id,
                   chat_id = :chat_id,
                   message_id = :message_id,
               	   update_ts = current_timestamp
               returning *;`
	entity := MessageEventEntityFromDomain(domain)
	saved, err := receiver.SaveEntity(query, entity)
	if err != nil || saved == nil {
		return nil, err
	}
	return saved.ToDomain(), nil
}
