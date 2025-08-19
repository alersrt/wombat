package daemon

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"wombat/pkg"
)

var ErrDaemon = errors.New("daemon")

const (
	ExitCodeDone         = 0
	ExitCodeError        = 1
	ExitCodeInvalidUsage = 2
	ExitCodeInterrupt    = 130
	ExitCodeKill         = 137
	ExitCodeTerminate    = 143
)

// HandleSignals handles os signals. Returns exit code and error if any.
func HandleSignals(ctx context.Context, cancel func()) (int, error) {
	slog.Info("daemon:handle:start")
	defer slog.Info("daemon:handle:finish")

	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan)

	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case s := <-signalChan:
			switch s {
			case os.Interrupt:
				cancel()
				slog.Info("daemon:handle:interrupt")
				return ExitCodeInterrupt, nil
			case os.Kill:
				slog.Info("daemon:handle:kill")
				return ExitCodeKill, nil
			case syscall.SIGTERM:
				cancel()
				slog.Info("daemon:handle:sigterm")
				return ExitCodeTerminate, nil
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				slog.Info("daemon:handle:ctx:err")
				return ExitCodeError, pkg.Wrap(ErrDaemon, err)
			} else {
				slog.Info("daemon:handle:ctx:done")
				return ExitCodeDone, nil
			}
		}
	}
}
