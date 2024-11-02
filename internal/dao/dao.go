package dao

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"wombat/internal/domain"
)

type QueryManager interface {
	SaveMessageEvent(ctx context.Context, entity *domain.MessageEvent) (result *domain.MessageEvent, err error)
}

type postgreSQLManager struct {
	url *string
	ctx context.Context
}

func NewPostgreSQLManager(ctx context.Context, url *string) QueryManager {
	return &postgreSQLManager{url: url, ctx: ctx}
}

func (receiver *postgreSQLManager) SaveMessageEvent(ctx context.Context, entity *domain.MessageEvent) (result *domain.MessageEvent, err error) {
	conn, err := sql.Open("pgx", *receiver.url)
	defer conn.Close()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	rows, err := conn.Query(
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

	saved := &domain.MessageEvent{}
	var sourceType, eventType string
	if rows.Next() {
		err = rows.Scan(&saved.Hash, &sourceType, &eventType, &saved.Text, &saved.AuthorId, &saved.ChatId, &saved.MessageId)
		err = saved.SourceType.FromString(sourceType)
		if err != nil {
			return nil, err
		}
		err = saved.EventType.FromString(eventType)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return saved, nil
}
