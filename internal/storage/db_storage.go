package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"wombat/internal/domain"
)

var ErrStorage = errors.New("storage")

type DbStorage struct {
	db *sqlx.DB
}

type Entity[D any] interface {
	ToDomain() *D
	FromDomain(*D) Entity[D]
}

func ToDomain[D any](entities []Entity[D]) []*D {
	var domains = make([]*D, len(entities))
	for i, entity := range entities {
		domains[i] = entity.ToDomain()
	}
	return domains
}

type Tx struct {
	tx *sqlx.Tx
}

func NewDbStorage(url string) (*DbStorage, error) {
	dbConn, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, errors.Join(ErrStorage, err)
	}
	return &DbStorage{db: dbConn}, nil
}

func (db *DbStorage) BeginTx(ctx context.Context) (*Tx, error) {
	if tx, err := db.db.BeginTxx(ctx, nil); err != nil {
		return nil, errors.Join(ErrStorage, err)
	} else {
		return &Tx{tx}, nil
	}
}

func (tx *Tx) CommitTx() error {
	if err := tx.tx.Commit(); err != nil {
		return errors.Join(ErrStorage, err)
	} else {
		return nil
	}
}

func (tx *Tx) RollbackTx() error {
	if err := tx.tx.Rollback(); err != nil {
		return errors.Join(ErrStorage, err)
	} else {
		return nil
	}
}

func (db *DbStorage) HasConnectionSource(sourceType string, userId string) (bool, error) {
	query := `select count(*)
              from wombatsm.source_connection wsc
              where wsc.source_type = $1
                and wsc.user_id = $2;`
	var count int
	if err := db.db.Get(&count, query, sourceType, userId); err != nil {
		return false, errors.Join(ErrStorage, err)
	} else {
		return count == 1, nil
	}
}

func (tx *Tx) CreateAccount() (*uuid.UUID, error) {
	gid := uuid.New()
	query := `insert into wombatsm.account(gid) values($1)`
	if _, err := tx.tx.Exec(query, &gid); err != nil {
		return nil, errors.Join(ErrStorage, err)
	} else {
		return &gid, nil
	}
}

func (tx *Tx) CreateSourceConnection(accountGid *uuid.UUID, sourceType string, userId string) error {
	query := `insert into wombatsm.source_connection(account_gid, source_type, user_id)
              values($1, $2, $3)`
	if _, err := tx.tx.Exec(query, accountGid, sourceType, userId); err != nil {
		return errors.Join(ErrStorage, err)
	} else {
		return nil
	}
}

func (tx *Tx) CreateTargetConnection(accountGid *uuid.UUID, targetType string, token []byte) error {
	query := `insert into wombatsm.target_connection(account_gid, target_type, token)
              values($1, $2, $3)`
	if _, err := tx.tx.Exec(query, accountGid, targetType, token); err != nil {
		return errors.Join(ErrStorage, err)
	} else {
		return nil
	}
}

func (tx *Tx) GetTargetConnection(sourceType string, targetType string, userId string) (*domain.TargetConnection, error) {
	query := `select wtc.*
              from wombatsm.account wa
                left join wombatsm.target_connection wtc on wa.gid = wtc.account_gid
                left join wombatsm.source_connection wsc on wa.gid = wsc.account_gid
              where wsc.source_type = $1
                and wtc.target_type = $2
                and wsc.user_id = $3`
	targetConnection := &TargetConnectionEntity{}
	if err := tx.tx.Get(targetConnection, query, sourceType, targetType, userId); err != nil {
		return nil, errors.Join(ErrStorage, err)
	} else {
		return targetConnection.ToDomain(), nil
	}
}

func (tx *Tx) GetCommentMetadata(sourceType string, chatId string, userId string) ([]*domain.Comment, error) {
	query := `select *
              from wombatsm.comment
              where source_type = $1
                and chat_id = $2
                and message_id = $3`
	rows, err := tx.tx.Queryx(query, sourceType, chatId, userId)
	if err != nil {
		return nil, errors.Join(ErrStorage, err)
	}
	var comments []Entity[domain.Comment]
	if rows.Next() {
		comment := &CommentEntity{}
		err = rows.Scan(comment)
		if err != nil {
			return nil, errors.Join(ErrStorage, err)
		}
		comments = append(comments, comment)
	}
	return ToDomain(comments), nil
}

func (tx *Tx) SaveCommentMetadata(domain *domain.Comment) (*domain.Comment, error) {
	query := `insert into wombatsm.comment(target_type, source_type, comment_id, user_id, chat_id, message_id, tag)
              values (:target_type, :source_type, :comment_id, :user_id, :chat_id, :message_id, :tag)
              returning *`
	rows, err := tx.tx.NamedQuery(query, (*CommentEntity).FromDomain(nil, domain))
	if err != nil {
		return nil, errors.Join(ErrStorage, err)
	}
	entity := &CommentEntity{}
	if rows.Next() {
		err = rows.Scan(entity)
		if err != nil {
			return nil, errors.Join(ErrStorage, err)
		}
	}
	return entity.ToDomain(), nil
}
