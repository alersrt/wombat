package pkg

import (
	"errors"
	"fmt"
	"runtime"
)

func Catch() {
	CatchWithReturn(nil)
}

func CatchWithReturn(err *error) {
	CatchWithReturnAndPost(err, nil)
}

func CatchWithReturnAndPost(err *error, post func()) {
	if r := recover(); r != nil {
		switch v := r.(type) {
		case LocatedError:
			fmt.Printf("%s:%d %v", v.filename, v.line, v.err)
			if err != nil {
				*err = v.err
			}
		case error:
			fmt.Printf("%v", v)
			if err != nil {
				*err = v
			}
		case string:
			fmt.Printf("%v", v)
			if err != nil {
				*err = errors.New(v)
			}
		default:
			fmt.Printf("%v", v)
		}

		if post != nil {
			post()
		}
	}
}

type LocatedError struct {
	filename string
	line     int
	err      error
}

func TryPanic(err error) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		panic(LocatedError{filename, line, err})
	}
}
