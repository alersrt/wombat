package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"wombat/pkg/log"
)

func Create(conf Config, task Task) {
	ctx, cancel := context.WithCancelCause(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		cancel(nil)
	}()

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					conf.Init(os.Args)
				case os.Interrupt:
					cancel(nil)
					os.Exit(1)
				}
			case <-ctx.Done():
				log.InfoLog.Print("Done.")
				os.Exit(1)
			}
		}
	}()

	if err := execute(ctx, cancel, conf, task); err != nil {
		log.ErrorLog.Printf("%s\n", err)
		os.Exit(1)
	}
}

func execute(ctx context.Context, cancel context.CancelCauseFunc, conf Config, task Task) error {
	conf.Init(os.Args)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			task(cancel)
		}
	}
}
