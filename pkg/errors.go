package pkg

import (
	_ "errors"
	"fmt"
	"strings"
)

// Wrap wraps errors in order of their declaration.
func Wrap(errors ...error) error {
	var sb strings.Builder
	if errors == nil || len(errors) == 0 {
		return nil
	}

	size := len(errors)
	args := make([]any, size)
	for i := 0; i < size; i++ {
		if i == size-1 {
			sb.WriteString("%w")
		} else {
			sb.WriteString("%w: ")
		}
		args[i] = errors[i]
	}

	return fmt.Errorf(sb.String(), args...)
}
