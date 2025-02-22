package dao

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"wombat/internal/domain"
)

type CommentRepository struct {
	*PostgreSQLQueryHelper[domain.Comment, uuid.UUID]
}

func NewCommentRepository(url *string) (*CommentRepository, error) {
	db, err := sqlx.Connect("postgres", *url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &CommentRepository{
		&PostgreSQLQueryHelper[domain.Comment, uuid.UUID]{
			db:            db,
			entityFactory: &CommentEntityFactory{},
		},
	}, nil
}

func (receiver *CommentRepository) GetById(id uuid.UUID) *domain.Comment {
	entity := receiver.GetEntityById("select * from wombatsm.comments where comment_id = $1", id)
	if entity == nil {
		return nil
	}
	return entity.ToDomain()
}

func (receiver *CommentRepository) Save(domain *domain.Comment) *domain.Comment {
	query := `insert into wombatsm.comments(
                              target_type,
                              source_type,
                              text,
                              author_id,
                              chat_id,
                              message_id,
                              comment_id,
                              tag)
               values (:target_type, :source_type, :text, :author_id, :chat_id, :message_id, :comment_id, :tag)
               on conflict (target_type, comment_id)
               do update
               set source_type = :source_type,
                   text = :text,
                   author_id = :author_id,
                   chat_id = :chat_id,
                   message_id = :message_id,
                   tag = :tag,
               	   update_ts = current_timestamp
               returning *;`
	entity := receiver.entityFactory.FromDomain(domain)
	saved := receiver.SaveEntity(query, entity)
	if saved == nil {
		return nil
	}
	return saved.ToDomain()
}

func (receiver *CommentRepository) GetMessagesByMetadata(sourceType domain.SourceType, chatId string, messageId string) []*domain.Comment {
	entities := receiver.GetEntitiesByArgs(`select *
                                                  from wombatsm.comments
                                                  where chat_id = $1
                                                    and message_id = $2
                                                    and source_type = $3`, chatId, messageId, sourceType)
	if entities == nil {
		return nil
	}

	var domains []*domain.Comment
	for _, entity := range entities {
		domains = append(domains, entity.ToDomain())
	}

	return domains
}
