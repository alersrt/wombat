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
	query := `insert into wombatsm.acl(author_id, is_allowed)
               values (:author_id, :is_allowed)
               on conflict (author_id)
               do update
               set author_id = :author_id,
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

func (receiver *AclRepository) IsAuthorAllowed(authorId string) bool {
	result := receiver.GetById(authorId)
	if result != nil {
		return result.IsAllowed
	} else {
		return false
	}
}
