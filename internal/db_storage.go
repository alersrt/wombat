package internal

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
)

type DbStorage struct {
	db *sqlx.DB
}

type Entity[D any] interface {
	ToDomain() *D
	FromDomain(*D) Entity[D]
}

func ToDomain[D any](entities []Entity[D]) []*D {
	var domains []*D
	for _, entity := range entities {
		domains = append(domains, entity.ToDomain())
	}
	return domains
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

func (receiver *Tx) CreateAccount() (*uuid.UUID, error) {
	gid := uuid.New()
	query := `insert into wombatsm.accounts(gid) values($1)`
	_, err := receiver.Exec(query, &gid)
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	} else {
		return &gid, nil
	}
}

func (receiver *Tx) CreateSourceConnection(accountGid *uuid.UUID, sourceType string, userId string) error {
	query := `insert into wombatsm.source_connections(account_gid, source_type, user_id)
              values($1, $2, $3)`
	_, err := receiver.Exec(query, accountGid, sourceType, userId)
	if err != nil {
		slog.Warn(err.Error())
		return err
	} else {
		return nil
	}
}

func (receiver *Tx) CreateTargetConnection(accountGid *uuid.UUID, targetType string, token string) error {
	query := `insert into wombatsm.target_connections(account_gid, target_type, token)
              values($1, $2, $3)`
	_, err := receiver.Exec(query, accountGid, targetType, token)
	if err != nil {
		slog.Warn(err.Error())
		return err
	} else {
		return nil
	}
}

func (receiver *Tx) GetTargetConnection(sourceType string, targetType string, userId string) (*TargetConnection, error) {
	query := `select wtc.*
              from wombatsm.accounts wa
                left join wombatsm.target_connections wtc on wa.gid = wtc.account_gid
                left join wombatsm.source_connections wsc on wa.gid = wsc.account_gid
              where wsc.source_type = $1
                and wtc.target_type = $2
                and wsc.user_id = $3`
	targetConnection := &TargetConnectionEntity{}
	err := receiver.Get(targetConnection, query, sourceType, targetType, userId)
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}
	return targetConnection.ToDomain(), nil
}

func (receiver *Tx) GetCommentMetadata(sourceType string, chatId string, userId string) ([]*Comment, error) {
	query := `select *
              from wombatsm.comments
              where source_type = $1
                and chat_id = $2
                and message_id = $3`
	rows, err := receiver.Queryx(query, sourceType, chatId, userId)
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}
	var comments []Entity[Comment]
	err = rows.Scan(comments)
	if err != nil {
		return nil, err
	}
	return ToDomain(comments), nil
}

func (receiver *Tx) SaveCommentMetadata(domain *Comment) (*Comment, error) {
	query := `insert into wombatsm.comments(target_type, source_type, comment_id, user_id, chat_id, message_id, tag)
              values (:target_type, :source_type, :comment_id, :user_id, :chat_id, :message_id, :tag)
              returning *`
	row := receiver.QueryRowx(query, (*CommentEntity).FromDomain(nil, domain))
	if row != nil && row.Err() != nil {
		slog.Warn(row.Err().Error())
		return nil, row.Err()
	}
	entity := &CommentEntity{}
	err := row.Scan(entity)
	if err != nil {
		return nil, err
	}
	return entity.ToDomain(), nil
}
