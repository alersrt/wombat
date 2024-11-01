package dao

import (
	"context"
	"github.com/jackc/pgx/v5"
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
	conn, err := pgx.Connect(ctx, *receiver.url)
	defer conn.Close(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	row, err := conn.Query(ctx,
		`insert into wombatsm.message_event(hash, source_type, event_type, text, author_id, chat_id, message_id)
               values (:hash,:sourceType,:eventType,:text,:authorId,:chatId,:messageId)
               on conflict (hash)
               do update
               set event_type = :eventType,
                   text = :text,
                   author_id = :authorId,
                   chat_id = :chatId,
                   message_id = :messageId
               returning hash, source_type, event_type, text, author_id, chat_id, message_id`,
		entity.Hash,
		entity.SourceType.String(),
		entity.EventType.String(),
		entity.Text,
		entity.AuthorId,
		entity.ChatId,
		entity.MessageId,
	)

	if err != nil {
		return nil, err
	}

	saved := &domain.MessageEvent{}
	if err := row.Scan(&saved.Hash, &saved.SourceType, &saved.EventType, &saved.Text, &saved.AuthorId, &saved.ChatId, &saved.MessageId); err != nil {
		return nil, err
	}
	return saved, nil
}
