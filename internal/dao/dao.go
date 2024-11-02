package dao

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"wombat/internal/domain"
)

type QueryHelper interface {
	GetMessageEvent(hash string) (*domain.MessageEvent, error)
	SaveMessageEvent(entity *domain.MessageEvent) (*domain.MessageEvent, error)
}

type PostgreSQLQueryHelper struct {
	db *sqlx.DB
}

func NewQueryHelper(url *string) (QueryHelper, error) {
	db, err := sqlx.Connect("postgres", *url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &PostgreSQLQueryHelper{db: db}, nil
}

func (receiver *PostgreSQLQueryHelper) GetMessageEvent(hash string) (*domain.MessageEvent, error) {
	entity := &MessageEventEntity{}
	err := receiver.db.QueryRowx("select * from wombatsm.message_event where hash = $1", hash).StructScan(*entity)
	if err != nil {
		return nil, err
	}
	return entity.ToDomain(), nil
}

func (receiver *PostgreSQLQueryHelper) SaveMessageEvent(domain *domain.MessageEvent) (*domain.MessageEvent, error) {
	entity := &MessageEventEntity{}
	entity.FromDomain(domain)

	result := &MessageEventEntity{}

	err := receiver.db.QueryRowx(
		`insert into wombatsm.message_event(hash, source_type, event_type, text, author_id, chat_id, message_id, create_ts, update_ts)
               values (:hash, :source_type, :event_type, :text, :author_id, :chat_id, :message_id, :create_ts, :update_ts)
               on conflict (hash)
               do update
               set event_type = :event_type,
                   text = :text,
                   author_id = :author_id,
                   chat_id = :chat_id,
                   message_id = :message_id,
               	   update_ts = :update_ts
               returning *;`,
		entity,
	).StructScan(*result)

	if err != nil {
		return nil, err
	}

	return result.ToDomain(), nil
}
