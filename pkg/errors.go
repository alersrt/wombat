package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
)

// Catch recovers the goroutine.
func Catch() {
	if r := recover(); r != nil {
	}
}

func CatchWithReturn(err *error) {
	if r := recover(); r != nil {
		switch v := r.(type) {
		case error:
			if err != nil {
				*err = v
			}
		case string:
			if err != nil {
				*err = errors.New(v)
			}
		default:
			if err != nil {
				*err = errors.New("internal")
			}
		}
	}
}

// Throw prints error and panic.
func Throw(err error) {
	if err != nil {
		err := errors.WithStack(err)
		slog.Error(fmt.Sprintf("%+v", err))
		panic(err)
	}
}
