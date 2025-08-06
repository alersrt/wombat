package internal

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"wombat/pkg"
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

func NewDbStorage(url string) (db *DbStorage, err error) {
	defer pkg.CatchWithReturn(&err)
	dbConn, err := sqlx.Connect("postgres", url)
	pkg.Throw(err)
	return &DbStorage{db: dbConn}, nil
}

func (db *DbStorage) BeginTx() *Tx {
	tx, err := db.db.Beginx()
	pkg.Throw(err)
	return &Tx{Tx: tx}
}

func (tx *Tx) CommitTx() {
	err := tx.Commit()
	pkg.Throw(err)
}

func (tx *Tx) RollbackTx() {
	err := tx.Rollback()
	pkg.Throw(err)
}

func (db *DbStorage) HasConnectionSource(sourceType string, userId string) bool {
	query := `select count(*)
              from wombatsm.source_connections wsc
              where wsc.source_type = $1
                and wsc.user_id = $2;`
	var count int
	err := db.db.Get(&count, query, sourceType, userId)
	pkg.Throw(err)
	return count == 1
}

func (tx *Tx) CreateAccount() *uuid.UUID {
	gid := uuid.New()
	query := `insert into wombatsm.accounts(gid) values($1)`
	_, err := tx.Exec(query, &gid)
	pkg.Throw(err)
	return &gid
}

func (tx *Tx) CreateSourceConnection(accountGid *uuid.UUID, sourceType string, userId string) {
	query := `insert into wombatsm.source_connections(account_gid, source_type, user_id)
              values($1, $2, $3)`
	_, err := tx.Exec(query, accountGid, sourceType, userId)
	pkg.Throw(err)
	return
}

func (tx *Tx) CreateTargetConnection(accountGid *uuid.UUID, targetType string, token []byte) {
	query := `insert into wombatsm.target_connections(account_gid, target_type, token)
              values($1, $2, $3, $4)`
	_, err := tx.Exec(query, accountGid, targetType, token)
	pkg.Throw(err)
	return
}

func (tx *Tx) GetTargetConnection(sourceType string, targetType string, userId string) *TargetConnection {
	query := `select wtc.*
              from wombatsm.accounts wa
                left join wombatsm.target_connections wtc on wa.gid = wtc.account_gid
                left join wombatsm.source_connections wsc on wa.gid = wsc.account_gid
              where wsc.source_type = $1
                and wtc.target_type = $2
                and wsc.user_id = $3`
	targetConnection := &TargetConnectionEntity{}
	err := tx.Get(targetConnection, query, sourceType, targetType, userId)
	pkg.Throw(err)
	return targetConnection.ToDomain()
}

func (tx *Tx) GetCommentMetadata(sourceType string, chatId string, userId string) []*Comment {
	query := `select *
              from wombatsm.comments
              where source_type = $1
                and chat_id = $2
                and message_id = $3`
	rows, err := tx.Queryx(query, sourceType, chatId, userId)
	pkg.Throw(err)
	var comments []Entity[Comment]
	if rows.Next() {
		comment := &CommentEntity{}
		err = rows.Scan(comment)
		comments = append(comments, comment)
	}
	pkg.Throw(err)
	return ToDomain(comments)
}

func (tx *Tx) SaveCommentMetadata(domain *Comment) *Comment {
	query := `insert into wombatsm.comments(target_type, source_type, comment_id, user_id, chat_id, message_id, tag)
              values (:target_type, :source_type, :comment_id, :user_id, :chat_id, :message_id, :tag)
              returning *`
	rows, err := tx.NamedQuery(query, (*CommentEntity).FromDomain(nil, domain))
	pkg.Throw(err)
	entity := &CommentEntity{}
	if rows.Next() {
		err = rows.Scan(entity)
	}
	pkg.Throw(err)
	return entity.ToDomain()
}
