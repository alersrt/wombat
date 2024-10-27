package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Create(conf Config, task Task) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	defer signal.Stop(signalChan)

	go handleSignals(ctx, cancel, conf, signalChan)

	if err := execute(ctx, cancel, conf, task); err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
}

func handleSignals(ctx context.Context, cancel context.CancelCauseFunc, conf Config, sigChan chan os.Signal) {
	for {
		select {
		case s := <-sigChan:
			switch s {
			case syscall.SIGHUP:
				conf.Init(os.Args)
			case os.Interrupt:
				cancel(nil)
				os.Exit(1)
			}
		case <-ctx.Done():
			slog.Info("Done.")
			os.Exit(1)
		}
	}
}

func execute(ctx context.Context, cancel context.CancelCauseFunc, conf Config, task Task) error {
	if !conf.IsInitiated() {
		err := conf.Init(os.Args)
		if err != nil {
			return err
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			task(cancel)
		}
	}
}
