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
func HandleSignals(ctx context.Context, cancel func()) (int, error) {
	slog.Info("daemon:handle:start")
	defer slog.Info("daemon:handle:finish")

	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	for {
		select {
		case s := <-signalChan:
			switch s {
			case os.Interrupt:
				cancel()
				slog.Info("daemon:handle:interrupt")
				return 130, nil
			case os.Kill:
				cancel()
				slog.Info("daemon:handle:kill")
				return 137, nil
			case syscall.SIGTERM:
				cancel()
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
