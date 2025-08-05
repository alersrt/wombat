package pkg

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

func Catch() {
	if r := recover(); r != nil {
		switch v := r.(type) {
		case LocatedError:
			fmt.Printf("%v %s:%d %v\n", v.createdTs, v.filename, v.line, v.err)
		case error:
			fmt.Printf("%v %v\n", time.Now(), v)
		case string:
			fmt.Printf("%v %v\n", time.Now(), v)
		default:
			fmt.Printf("%v %v\n", time.Now(), v)
		}
	}
}

func CatchWithReturn(err *error) {
	if r := recover(); r != nil {
		switch v := r.(type) {
		case LocatedError:
			fmt.Printf("%v %s:%d %v\n", v.createdTs, v.filename, v.line, v.err)
			if err != nil {
				*err = v.err
			}
		case error:
			fmt.Printf("%v %v\n", time.Now(), v)
			if err != nil {
				*err = v
			}
		case string:
			fmt.Printf("%v %v\n", time.Now(), v)
			if err != nil {
				*err = errors.New(v)
			}
		default:
			fmt.Printf("%v %v\n", time.Now(), v)
		}
	}
}

func CatchWithReturnAndPost(err *error, post func()) {
	if r := recover(); r != nil {
		switch v := r.(type) {
		case LocatedError:
			fmt.Printf("%v %s:%d %v\n", v.createdTs, v.filename, v.line, v.err)
			if err != nil {
				*err = v.err
			}
		case error:
			fmt.Printf("%v %v\n", time.Now(), v)
			if err != nil {
				*err = v
			}
		case string:
			fmt.Printf("%v %v\n", time.Now(), v)
			if err != nil {
				*err = errors.New(v)
			}
		default:
			fmt.Printf("%v %v\n", time.Now(), v)
		}

		if post != nil {
			post()
		}
	}
}

type LocatedError struct {
	createdTs time.Time
	filename  string
	line      int
	err       error
}

func (l *LocatedError) Error() string {
	return l.err.Error()
}

func TryPanic(err error) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		panic(LocatedError{time.Now(), filename, line, err})
	}
}
