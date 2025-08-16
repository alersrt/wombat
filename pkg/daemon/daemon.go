package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// HandleSignals handles os signals. Returns exit code and error if any.
func HandleSignals(ctx context.Context, shutdown func(), init func() error) (int, error) {
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
				shutdown()
				if err := init(); err != nil {
					slog.Info("daemon:handle:sighup:err")
					return 1, fmt.Errorf("daemon:handle:sighup:err: %w", err)
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
				slog.Info("daemon:handle:ctx:err")
				return 1, fmt.Errorf("daemon:handle:ctx:err: %w", err)
			} else {
				slog.Info("daemon:handle:ctx:done")
				return 0, nil
			}
		}
	}
}
