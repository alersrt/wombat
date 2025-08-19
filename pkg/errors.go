package pkg

import (
	_ "errors"
	"fmt"
	"strings"
)

func Wrap(errors ...any) error {
	var sb strings.Builder
	if errors == nil || len(errors) == 0 {
		return nil
	}

	size := len(errors)
	for i := 0; i < size; i++ {
		if i == size-1 {
			sb.WriteString("%w")
		} else {
			sb.WriteString("%w: ")
		}
	}
	return fmt.Errorf(sb.String(), errors...)
}
