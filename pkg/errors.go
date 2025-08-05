package pkg

import (
	"errors"
	"fmt"
	"log/slog"
)

func Catch() {
	if r := recover(); r != nil {
		switch v := r.(type) {
		case error:
			slog.Error(fmt.Sprintf("%+v", v))
		case string:
			slog.Error(fmt.Sprintf("%s", v))
		default:
			slog.Error(fmt.Sprintf("%+v", v))
		}
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
		}
	}
}

func CatchWithReturnAndPost(err *error, post func()) {
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
		}

		if post != nil {
			post()
		}
	}
}

func TryPanic(err error) {
	if err != nil {
		panic(err)
	}
}
