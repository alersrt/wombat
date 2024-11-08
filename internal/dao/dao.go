package dao

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
)

type Entity[D any] interface {
	ToDomain() *D
}

type QueryHelper[D any, ID any] interface {
	GetById(id ID) *D
	Save(domain *D) *D
}

type EntityFactory[D any] interface {
	EmptyEntity() Entity[D]
	FromDomain(*D) Entity[D]
}

type PostgreSQLQueryHelper[D any, ID any] struct {
	db            *sqlx.DB
	entityFactory EntityFactory[D]
}

func (receiver *PostgreSQLQueryHelper[D, ID]) GetEntityById(query string, id ID) Entity[D] {
	entity := receiver.entityFactory.EmptyEntity()
	err := receiver.db.QueryRowx(query, id).StructScan(entity)
	if err != nil {
		slog.Warn(err.Error())
		return nil
	}
	return entity
}

func (receiver *PostgreSQLQueryHelper[D, ID]) GetEntitiesByArgs(query string, args ...any) []Entity[D] {
	rows, err := receiver.db.Queryx(query, args...)
	if err != nil {
		slog.Warn(err.Error())
		return nil
	}

	var entities []Entity[D]
	if rows.Next() {
		entity := receiver.entityFactory.EmptyEntity()
		err = rows.StructScan(entity)
		if err != nil {
			slog.Warn(err.Error())
			return nil
		}
		entities = append(entities, entity)
	}

	return entities
}

func (receiver *PostgreSQLQueryHelper[D, ID]) SaveEntity(query string, entity Entity[D]) Entity[D] {
	rows, err := receiver.db.NamedQuery(query, entity)
	if err != nil {
		slog.Warn(err.Error())
		return nil
	}

	if rows.Next() {
		saved := receiver.entityFactory.EmptyEntity()
		err = rows.StructScan(saved)
		if err != nil {
			slog.Warn(err.Error())
			return nil
		}
		return saved
	}

	return nil
}
