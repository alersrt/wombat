package dao

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5"
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
	conn, err := sql.Open("postgres", *receiver.url)
	defer conn.Close()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	row := conn.QueryRow(
		`insert into wombatsm.message_event(hash, source_type, event_type, text, author_id, chat_id, message_id)
               values (@hash,@sourceType,@eventType,@text,@authorId,@chatId,@messageId)
               on conflict (hash)
               do update
               set event_type = @eventType,
                   text = @text,
                   author_id = @authorId,
                   chat_id = @chatId,
                   message_id = @messageId
               returning hash, source_type, event_type, text, author_id, chat_id, message_id`,
		sql.Named("hash", entity.Hash),
		sql.Named("sourceType", entity.SourceType.String()),
		sql.Named("eventType", entity.EventType.String()),
		sql.Named("text", entity.Text),
		sql.Named("authorId", entity.AuthorId),
		sql.Named("chatId", entity.ChatId),
		sql.Named("messageUd", entity.MessageId),
	)

	if err != nil {
		return nil, err
	}

	saved := &domain.MessageEvent{}
	err = row.Scan(&saved.Hash, &saved.SourceType, &saved.EventType, &saved.Text, &saved.AuthorId, &saved.ChatId, &saved.MessageId)
	if err != nil {
		return nil, err
	}
	return saved, nil
}
