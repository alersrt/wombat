package app

import (
	"wombat/internal/domain"
)

type Source interface {
	GetSourceType() domain.SourceType
	Process()
}

type Target interface {
	GetTargetType() domain.TargetType
	Process()
}
