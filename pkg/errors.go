package pkg

import "fmt"

func Catch(err error) {
	if r := recover(); r != nil {
		var ok bool
		err, ok = r.(error)
		if !ok {
			err = fmt.Errorf("pkg: %v", r)
		}
	}
}

func CatchWithoutReturn() {
	if r := recover(); r != nil {
		fmt.Printf("pkg: %v", r)
	}
}

func CatchWithPost(err error, post func()) {
	if r := recover(); r != nil {
		var ok bool
		err, ok = r.(error)
		if !ok {
			err = fmt.Errorf("pkg: %v", r)
		}

		post()
	}
}

func Try(err error) {
	if err != nil {
		panic(err)
	}
}
