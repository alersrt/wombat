package dao

import "wombat/internal/domain"

type AclEntity struct {
	AuthorId  string `db:"author_id"`
	IsAllowed bool   `db:"is_allowed"`
}

func (receiver *AclEntity) ToDomain() *domain.Acl {
	return &domain.Acl{
		AuthorId:  receiver.AuthorId,
		IsAllowed: receiver.IsAllowed,
	}
}

type AclEntityFactory struct{}

func (receiver *AclEntityFactory) EmptyEntity() Entity[domain.Acl] {
	return &AclEntity{}
}

func (receiver *AclEntityFactory) FromDomain(domain *domain.Acl) Entity[domain.Acl] {
	return &AclEntity{
		AuthorId:  domain.AuthorId,
		IsAllowed: domain.IsAllowed,
	}
}
