package daemon

import (
	"context"
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

	if err := execute(cancel, conf, task); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func handleSignals(ctx context.Context, cancel context.CancelCauseFunc, conf Config, sigChan chan os.Signal) {
	for {
		select {
		case s := <-sigChan:
			switch s {
			case syscall.SIGHUP:
				err := conf.Init(os.Args)
				if err != nil {
					slog.Error(err.Error())
				}
			case os.Interrupt:
				cancel(nil)
				os.Exit(1)
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				slog.Error(err.Error())
			}
			slog.Info("Done.")
			os.Exit(1)
		}
	}
}

func execute(cancel context.CancelCauseFunc, conf Config, task Task) error {
	if !conf.IsInitiated() {
		err := conf.Init(os.Args)
		if err != nil {
			return err
		}
	}

	for {
		task(cancel)
	}
}
