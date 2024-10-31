package daemon

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Daemon struct {
	ctx             context.Context
	cancelCauseFunc context.CancelCauseFunc
	conf            Config
}

func Create(ctx context.Context, cancelCauseFunc context.CancelCauseFunc, conf Config) *Daemon {
	dmn := &Daemon{
		ctx:             ctx,
		cancelCauseFunc: cancelCauseFunc,
		conf:            conf,
	}

	go dmn.handleSignals()

	return dmn
}

func (receiver *Daemon) Start(task Task) {
	if !receiver.conf.IsInitiated() {
		err := receiver.conf.Init(os.Args)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}

	for {
		task()
	}
}

func (receiver *Daemon) GetConfig() Config {
	return receiver.conf
}
func (receiver *Daemon) GetContext() context.Context {
	return receiver.ctx
}
func (receiver *Daemon) GetCancelCauseFunc() context.CancelCauseFunc {
	return receiver.cancelCauseFunc
}

func (receiver *Daemon) handleSignals() {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	defer signal.Stop(signalChan)

	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP:
			err := receiver.conf.Init(os.Args)
			if err != nil {
				slog.Error(err.Error())
			}
		case os.Interrupt:
			receiver.cancelCauseFunc(nil)
			os.Exit(130)
		case os.Kill:
			os.Exit(137)
		case syscall.SIGTERM:
			receiver.cancelCauseFunc(nil)
			os.Exit(143)
		}
	case <-receiver.ctx.Done():
		if err := receiver.ctx.Err(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		} else {
			slog.Info("Done.")
			os.Exit(0)
		}
	}
}
