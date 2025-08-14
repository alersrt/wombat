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

// HandleSignals handles os signals. Returns exit code and error if any.
func HandleSignals(ctx context.Context, cfg Config) (int, error) {
	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	for {
		select {
		case s := <-signalChan:
			switch s {
			case syscall.SIGHUP:
				if err := cfg.Init(os.Args); err != nil {
					return 1, err
				}
			case os.Interrupt:
				return 130, nil
			case os.Kill:
				return 137, nil
			case syscall.SIGTERM:
				return 143, nil
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return 1, err
			} else {
				slog.Info("Done.")
				return 0, nil
			}
		}
	}
}
