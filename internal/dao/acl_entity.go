package dao

import (
	"time"
	"wombat/internal/domain"
)

type AclEntity struct {
	SourceType domain.SourceType `db:"source_type"`
	AuthorId   string            `db:"author_id"`
	IsAllowed  bool              `db:"is_allowed"`
	CreateTs   time.Time         `db:"create_ts"`
	UpdateTs   time.Time         `db:"update_ts"`
}

func (receiver *AclEntity) ToDomain() *domain.Acl {
	return &domain.Acl{
		SourceType: receiver.SourceType,
		AuthorId:   receiver.AuthorId,
		IsAllowed:  receiver.IsAllowed,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

type AclEntityFactory struct{}

func (receiver *AclEntityFactory) EmptyEntity() Entity[domain.Acl] {
	return &AclEntity{}
}

func (receiver *AclEntityFactory) FromDomain(domain *domain.Acl) Entity[domain.Acl] {
	return &AclEntity{
		SourceType: domain.SourceType,
		AuthorId:   domain.AuthorId,
		IsAllowed:  domain.IsAllowed,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}
