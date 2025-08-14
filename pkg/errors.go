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

// Throw prints error and panic.
func Throw(err error) {
	if err != nil {
		err := errors.WithStack(err)
		slog.Error(fmt.Sprintf("%+v", err))
		panic(err)
	}
}
