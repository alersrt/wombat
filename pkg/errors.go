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
			if err != nil {
				*err = errors.New("internal")
			}
		}
	}
}

func CatchWithReturnAndCall(err *error, call func()) {
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

		if call != nil {
			call()
		}
	}
}

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}
