package dao

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
	"wombat/internal/domain"
)

type MessageEventRepository struct {
	*PostgreSQLQueryHelper[domain.MessageEvent, string]
}

func NewMessageEventRepository(url *string) (*MessageEventRepository, error) {
	db, err := sqlx.Connect("postgres", *url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &MessageEventRepository{
		&PostgreSQLQueryHelper[domain.MessageEvent, string]{
			db:            db,
			entityFactory: &MessageEventEntityFactory{},
		},
	}, nil
}

func (receiver *MessageEventRepository) GetById(id string) *domain.MessageEvent {
	entity := receiver.GetEntityById("select * from wombatsm.message_event where hash = $1", id)
	if entity == nil {
		return nil
	}
	return entity.ToDomain()
}

func (receiver *MessageEventRepository) Save(domain *domain.MessageEvent) *domain.MessageEvent {
	query := `insert into wombatsm.message_event(hash, source_type, text, author_id, chat_id, message_id, comment_id)
               values (:hash, :source_type, :text, :author_id, :chat_id, :message_id, :comment_id)
               on conflict (hash)
               do update
               set hash = :hash,
                   source_type = :source_type,
                   text = :text,
                   author_id = :author_id,
                   chat_id = :chat_id,
                   message_id = :message_id,
                   comment_id = :comment_id,
               	   update_ts = current_timestamp
               returning *;`
	entity := receiver.entityFactory.FromDomain(domain)
	saved := receiver.SaveEntity(query, entity)
	if saved == nil {
		return nil
	}
	return saved.ToDomain()
}
