package dao

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"wombat/internal/domain"
)

type AclRepository struct {
	*PostgreSQLQueryHelper[domain.Acl, uuid.UUID]
}

func NewAclRepository(url *string) (*AclRepository, error) {
	db, err := sqlx.Connect("postgres", *url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &AclRepository{
		&PostgreSQLQueryHelper[domain.Acl, uuid.UUID]{
			db:            db,
			entityFactory: &AclEntityFactory{},
		},
	}, nil
}

func (receiver *AclRepository) GetById(id uuid.UUID) *domain.Acl {
	entity := receiver.GetEntityById("select * from wombatsm.acl where gid = $1", id)
	if entity == nil {
		return nil
	}
	return entity.ToDomain()
}

func (receiver *AclRepository) Save(domain *domain.Acl) *domain.Acl {
	query := `insert into wombatsm.acl(source_type, author_id, is_allowed)
               values (:source_type, :author_id, :is_allowed)
               on conflict (source_type, author_id)
               do update
               set source_type = :source_type,
                   author_id = :author_id,
                   is_allowed = :is_allowed,
               	   update_ts = current_timestamp
               returning *;`
	entity := receiver.entityFactory.FromDomain(domain)
	saved := receiver.SaveEntity(query, entity)
	if saved == nil {
		return nil
	}
	return saved.ToDomain()
}

func (receiver *AclRepository) IsAuthorAllowed(sourceType domain.SourceType, authorId string) bool {
	query := `select *
              from wombatsm.acl
              where source_type = $1
                and author_id = $2;`
	result := receiver.GetEntitiesByArgs(query, sourceType, authorId)
	if result != nil && len(result) == 1 {
		return result[0].ToDomain().IsAllowed
	} else {
		return false
	}
}
