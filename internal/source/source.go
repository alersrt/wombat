package source

import "context"

type Source interface {
	Read(causeFunc context.CancelCauseFunc)
}
