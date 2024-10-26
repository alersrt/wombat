package daemon

type Config interface {
	Init(args []string) error
	IsInitiated() bool
}
