package pkg

type Args[T any] struct {
	Headers map[string]string `json:"headers"`
	Value   T                 `json:"value"`
}
