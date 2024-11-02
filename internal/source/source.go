package source

type Source interface {
	ForwardTo(chan any)
}
