package dao

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"wombat/internal/domain"
)

type ConnectionRepository struct {
	*PostgreSQLQueryHelper[domain.Connection, uuid.UUID]
}

func NewConnectionRepository(url *string) (*ConnectionRepository, error) {
	db, err := sqlx.Connect("postgres", *url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &ConnectionRepository{
		&PostgreSQLQueryHelper[domain.Connection, uuid.UUID]{
			db:            db,
			entityFactory: &ConnectionEntityFactory{},
		},
	}, nil
}

func (receiver *ConnectionRepository) GetById(id uuid.UUID) *domain.Connection {
	entity := receiver.GetEntityById("select * from wombatsm.connections where gid = $1", id)
	if entity == nil {
		return nil
	}
	return entity.ToDomain()
}

func (receiver *ConnectionRepository) Save(domain *domain.Connection) *domain.Connection {
	query := `insert into wombatsm.connections(target_type, token, source_type, author_id)
               values (:target_type, :token, :source_type, :author_id)
               on conflict (target_type, source_type, author_id)
               do update
               set token = :token,
               	   update_ts = current_timestamp
               returning *;`
	entity := receiver.entityFactory.FromDomain(domain)
	saved := receiver.SaveEntity(query, entity)
	if saved == nil {
		return nil
	}
	return saved.ToDomain()
}

func (receiver *ConnectionRepository) GetToken(targetType domain.TargetType, sourceType domain.SourceType, authorId string) (url string, token string) {
	query := `select *
              from wombatsm.connections
              where target_type = $1
                and source_type = $2
                and author_id = $3;`
	result := receiver.GetEntitiesByArgs(query, targetType.String(), sourceType.String(), authorId)
	if result != nil && len(result) == 1 {
		return "", result[0].ToDomain().Token
	} else {
		return "", ""
	}
}
