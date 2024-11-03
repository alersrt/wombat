package dao

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Entity[D any] interface {
	ToDomain() *D
}

type QueryHelper[D any, ID any] interface {
	GetById(id ID) (*D, error)
	Save(domain *D) (*D, error)
}

type EntityFactory[D any] interface {
	EmptyEntity() Entity[D]
	FromDomain(*D) Entity[D]
}

type PostgreSQLQueryHelper[D any, ID any] struct {
	db            *sqlx.DB
	entityFactory EntityFactory[D]
}

func (receiver *PostgreSQLQueryHelper[D, ID]) GetEntityById(query string, id ID) (Entity[D], error) {
	entity := receiver.entityFactory.EmptyEntity()
	err := receiver.db.QueryRowx(query, id).StructScan(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (receiver *PostgreSQLQueryHelper[D, ID]) SaveEntity(query string, entity Entity[D]) (Entity[D], error) {
	rows, err := receiver.db.NamedQuery(query, entity)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		saved := receiver.entityFactory.EmptyEntity()
		err = rows.StructScan(saved)
		if err != nil {
			return nil, err
		}
		return saved, nil
	}

	return nil, nil
}
