package daemon

import "context"

type Task func(context.CancelCauseFunc)
