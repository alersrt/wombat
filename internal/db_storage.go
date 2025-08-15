package internal

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

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
		return nil, errors.New(err.Error())
	}
	return &DbStorage{db: dbConn}, nil
}

func (db *DbStorage) BeginTx(ctx context.Context) (*Tx, error) {
	if tx, err := db.db.BeginTxx(ctx, nil); err != nil {
		return nil, errors.New(err.Error())
	} else {
		return &Tx{tx}, nil
	}
}

func (tx *Tx) CommitTx() error {
	if err := tx.tx.Commit(); err != nil {
		return errors.New(err.Error())
	} else {
		return nil
	}
}

func (tx *Tx) RollbackTx() error {
	if err := tx.tx.Rollback(); err != nil {
		return errors.New(err.Error())
	} else {
		return nil
	}
}

func (db *DbStorage) HasConnectionSource(sourceType string, userId string) (bool, error) {
	query := `select count(*)
              from wombatsm.source_connections wsc
              where wsc.source_type = $1
                and wsc.user_id = $2;`
	var count int
	if err := db.db.Get(&count, query, sourceType, userId); err != nil {
		return false, errors.New(err.Error())
	} else {
		return count == 1, nil
	}
}

func (tx *Tx) CreateAccount() (*uuid.UUID, error) {
	gid := uuid.New()
	query := `insert into wombatsm.accounts(gid) values($1)`
	if _, err := tx.tx.Exec(query, &gid); err != nil {
		return nil, errors.New(err.Error())
	} else {
		return &gid, nil
	}
}

func (tx *Tx) CreateSourceConnection(accountGid *uuid.UUID, sourceType string, userId string) error {
	query := `insert into wombatsm.source_connections(account_gid, source_type, user_id)
              values($1, $2, $3)`
	if _, err := tx.tx.Exec(query, accountGid, sourceType, userId); err != nil {
		return errors.New(err.Error())
	} else {
		return nil
	}
}

func (tx *Tx) CreateTargetConnection(accountGid *uuid.UUID, targetType string, token []byte) error {
	query := `insert into wombatsm.target_connections(account_gid, target_type, token)
              values($1, $2, $3)`
	if _, err := tx.tx.Exec(query, accountGid, targetType, token); err != nil {
		return errors.New(err.Error())
	} else {
		return nil
	}
}

func (tx *Tx) GetTargetConnection(sourceType string, targetType string, userId string) (*TargetConnection, error) {
	query := `select wtc.*
              from wombatsm.accounts wa
                left join wombatsm.target_connections wtc on wa.gid = wtc.account_gid
                left join wombatsm.source_connections wsc on wa.gid = wsc.account_gid
              where wsc.source_type = $1
                and wtc.target_type = $2
                and wsc.user_id = $3`
	targetConnection := &TargetConnectionEntity{}
	if err := tx.tx.Get(targetConnection, query, sourceType, targetType, userId); err != nil {
		return nil, errors.New(err.Error())
	} else {
		return targetConnection.ToDomain(), nil
	}
}

func (tx *Tx) GetCommentMetadata(sourceType string, chatId string, userId string) ([]*Comment, error) {
	query := `select *
              from wombatsm.comments
              where source_type = $1
                and chat_id = $2
                and message_id = $3`
	rows, err := tx.tx.Queryx(query, sourceType, chatId, userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	var comments []Entity[Comment]
	if rows.Next() {
		comment := &CommentEntity{}
		err = rows.Scan(comment)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		comments = append(comments, comment)
	}
	return ToDomain(comments), nil
}

func (tx *Tx) SaveCommentMetadata(domain *Comment) (*Comment, error) {
	query := `insert into wombatsm.comments(target_type, source_type, comment_id, user_id, chat_id, message_id, tag)
              values (:target_type, :source_type, :comment_id, :user_id, :chat_id, :message_id, :tag)
              returning *`
	rows, err := tx.tx.NamedQuery(query, (*CommentEntity).FromDomain(nil, domain))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	entity := &CommentEntity{}
	if rows.Next() {
		err = rows.Scan(entity)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	}
	return entity.ToDomain(), nil
}
