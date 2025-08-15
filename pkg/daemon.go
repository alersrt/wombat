package pkg

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Task interface {
	Do(ctx context.Context)
}

type Config interface {
	Init(args []string) error
}

type Daemon interface {
	Init(args []string) error
	Shutdown()
}

// HandleSignals handles os signals. Returns exit code and error if any.
func HandleSignals(ctx context.Context, daemon Daemon) (int, error) {
	slog.Info("daemon:handle:start")
	defer slog.Info("daemon:handle:finish")

	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	for {
		select {
		case s := <-signalChan:
			switch s {
			case syscall.SIGHUP:
				slog.Info("daemon:handle:sighup:start")
				daemon.Shutdown()
				if err := daemon.Init(os.Args); err != nil {
					slog.Info("daemon:handle:sighup:error")
					return 1, err
				}
				slog.Info("daemon:handle:sighup:finish")
			case os.Interrupt:
				slog.Info("daemon:handle:interrupt")
				return 130, nil
			case os.Kill:
				slog.Info("daemon:handle:kill")
				return 137, nil
			case syscall.SIGTERM:
				slog.Info("daemon:handle:sigterm")
				return 143, nil
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				slog.Info("daemon:handle:ctx:error")
				return 1, err
			} else {
				slog.Info("daemon:handle:ctx:done")
				return 0, nil
			}
		}
	}
}
