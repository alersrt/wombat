package storage

import (
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

type Tx struct {
	*sqlx.Tx
}

func NewDbStorage(url string) (*DbStorage, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &DbStorage{
		db: db,
	}, nil
}

func (receiver *DbStorage) BeginTx() (*Tx, error) {
	tx, err := receiver.db.Beginx()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &Tx{
		Tx: tx,
	}, nil
}

func (receiver *Tx) CommitTx() error {
	return receiver.Commit()
}

func (receiver *Tx) RollbackTx() error {
	return receiver.Rollback()
}

func (receiver *DbStorage) HasConnectionSource(sourceType string, userId string) bool {
	query := `select count(*)
              from wombatsm.source_connections wsc
              where wsc.source_type = $1
                and wsc.user_id = $2;`
	var count int
	err := receiver.db.Get(&count, query, sourceType, userId)
	if err != nil {
		slog.Warn(err.Error())
		return false
	}
	return count == 1
}

func (receiver *Tx) SaveComment(domain *domain.Comment) error {
	query := `insert into wombatsm.comments(
                              target_type,
                              source_type,
                              comment_id,
                              user_id,
                              chat_id,
                              message_id,
                              tag)
               values (:target_type, :source_type, :user_id, :chat_id, :message_id, :comment_id, :tag)
               on conflict (target_type, comment_id)
               do update
               set source_type = :source_type,
                   user_id = :user_id,
                   chat_id = :chat_id,
                   message_id = :message_id,
                   tag = :tag,
               	   update_ts = current_timestamp`
	row := receiver.QueryRowx(query, (*CommentEntity).FromDomain(nil, domain))
	return row.Err()
}

func (receiver *Tx) GetCommentsByMetadata(sourceType domain.SourceType, chatId string, messageId string) ([]*domain.Comment, error) {
	query := `select *
               from wombatsm.comments
               where chat_id = $1
                 and message_id = $2
                 and source_type = $3`

	entities := []CommentEntity{}
	err := receiver.Select(entities, query, chatId, messageId, sourceType.String())

	if err != nil {
		return nil, err
	}

	var domains []*domain.Comment
	for _, entity := range entities {
		domains = append(domains, entity.ToDomain())
	}

	return domains, nil
}
