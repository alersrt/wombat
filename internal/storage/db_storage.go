package storage

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"wombat/internal/domain"
)

type DbStorage struct {
	db *sqlx.DB
}

type Entity[D any] interface {
	ToDomain() *D
	FromDomain(*D) Entity[D]
}

func NewDbStorage(url *string) (*DbStorage, error) {
	db, err := sqlx.Connect("postgres", *url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &DbStorage{
		db: db,
	}, nil
}

func (receiver *DbStorage) GetAclById(id uuid.UUID) (*domain.Acl, error) {
	entity := &AclEntity{}
	err := receiver.db.Get("select * from wombatsm.acl where gid = $1", id.String())
	if err != nil {
		return nil, err
	} else {
		return entity.ToDomain(), nil
	}
}

func (receiver *DbStorage) SaveAcl(domain *domain.Acl) (*domain.Acl, error) {
	query := `insert into wombatsm.acl(source_type, author_id, is_allowed)
               values (:source_type, :author_id, :is_allowed)
               on conflict (source_type, author_id)
               do update
               set source_type = :source_type,
                   author_id = :author_id,
                   is_allowed = :is_allowed,
               	   update_ts = current_timestamp
               returning *;`
	row := receiver.db.QueryRowx(query, (*AclEntity).FromDomain(nil, domain))

	if row.Err() != nil {
		return nil, row.Err()
	}

	entity := &AclEntity{}
	err := row.StructScan(entity)
	if err != nil {
		return nil, err
	}
	return entity.ToDomain(), nil
}

func (receiver *DbStorage) IsAuthorAllowed(sourceType domain.SourceType, authorId string) (bool, error) {
	query := `select *
              from wombatsm.acl
              where source_type = $1
                and author_id = $2`

	entities := []AclEntity{}
	err := receiver.db.Select(entities, query, sourceType.String(), authorId)

	if err != nil {
		return false, err
	} else {
		if entities != nil && len(entities) == 1 {
			return entities[0].ToDomain().IsAllowed, nil
		} else {
			return false, nil
		}
	}

}

func (receiver *DbStorage) GetCommentById(id uuid.UUID) (*domain.Comment, error) {
	entity := &CommentEntity{}
	err := receiver.db.Get(&entity, "select * from wombatsm.comments where comment_id = $1", id)
	if err != nil {
		return nil, err
	} else {
		return entity.ToDomain(), nil
	}
}

func (receiver *DbStorage) SaveComment(domain *domain.Comment) (*domain.Comment, error) {
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
	row := receiver.db.QueryRowx(query, (*CommentEntity).FromDomain(nil, domain))

	if row.Err() != nil {
		return nil, row.Err()
	}

	entity := &CommentEntity{}
	err := row.StructScan(entity)
	if err != nil {
		return nil, err
	}
	return entity.ToDomain(), nil
}

func (receiver *DbStorage) GetCommentsByMetadata(sourceType domain.SourceType, chatId string, messageId string) ([]*domain.Comment, error) {
	query := `select *
               from wombatsm.comments
               where chat_id = $1
                 and message_id = $2
                 and source_type = $3`

	entities := []CommentEntity{}
	err := receiver.db.Select(entities, query, chatId, messageId, sourceType.String())

	if err != nil {
		return nil, err
	}

	var domains []*domain.Comment
	for _, entity := range entities {
		domains = append(domains, entity.ToDomain())
	}

	return domains, nil
}
