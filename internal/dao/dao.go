package dao

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"wombat/internal/domain"
)

type QueryHelper interface {
	GetMessageEvent(ctx context.Context, hash string) (*domain.MessageEvent, error)
	SaveMessageEvent(ctx context.Context, entity *domain.MessageEvent) (*domain.MessageEvent, error)
}

type postgreSQLQueryHelper struct {
	db  *sql.DB
	ctx context.Context
}

func NewQueryHelper(ctx context.Context, url *string) (QueryHelper, error) {
	db, err := sql.Open("pgx", *url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &postgreSQLQueryHelper{db: db, ctx: ctx}, nil
}

func (receiver *postgreSQLQueryHelper) GetMessageEvent(ctx context.Context, hash string) (*domain.MessageEvent, error) {
	rows, err := receiver.db.Query(
		`select hash, source_type, event_type, text, author_id, chat_id, message_id
               from wombatsm.message_event
               where hash = $1`,
		hash,
	)
	if err != nil {
		return nil, err
	}

	res, err := scanMessageEvents(ctx, rows)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func (receiver *postgreSQLQueryHelper) SaveMessageEvent(ctx context.Context, entity *domain.MessageEvent) (*domain.MessageEvent, error) {
	rows, err := receiver.db.Query(
		`insert into wombatsm.message_event(hash, source_type, event_type, text, author_id, chat_id, message_id)
               values (@hashParam, @sourceType, @eventType, @textParam, @authorId, @chatId, @messageId)
               on conflict (hash)
               do update
               set event_type = @eventType,
                   text = @textParam,
                   author_id = @authorId,
                   chat_id = @chatId,
                   message_id = @messageId
               returning hash, source_type, event_type, text, author_id, chat_id, message_id;`,
		pgx.NamedArgs{
			"hashParam":  entity.Hash,
			"sourceType": entity.SourceType.String(),
			"eventType":  entity.EventType.String(),
			"textParam":  entity.Text,
			"authorId":   entity.AuthorId,
			"chatId":     entity.ChatId,
			"messageId":  entity.MessageId,
		})
	if err != nil {
		return nil, err
	}

	res, err := scanMessageEvents(ctx, rows)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func scanMessageEvent(ctx context.Context, row *sql.Row) (*domain.MessageEvent, error) {
	saved := &domain.MessageEvent{}
	var sourceType, eventType string

	err := row.Scan(&saved.Hash, &sourceType, &eventType, &saved.Text, &saved.AuthorId, &saved.ChatId, &saved.MessageId)

	if err != nil {
		return nil, err
	}

	saved.SourceType.FromString(sourceType)
	saved.EventType.FromString(eventType)

	return saved, nil
}

func scanMessageEvents(ctx context.Context, rows *sql.Rows) ([]*domain.MessageEvent, error) {
	var results []*domain.MessageEvent
	for rows.Next() {
		saved := &domain.MessageEvent{}
		var sourceType, eventType string
		err := rows.Scan(&saved.Hash, &sourceType, &eventType, &saved.Text, &saved.AuthorId, &saved.ChatId, &saved.MessageId)
		if err != nil {
			return nil, err
		}
		saved.SourceType.FromString(sourceType)
		saved.EventType.FromString(eventType)

		results = append(results, saved)
	}

	return results, nil
}
