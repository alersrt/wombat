package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	. "wombat/pkg/log"
)

func Create(conf Config, task func()) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					conf.Init(os.Args)
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				InfoLog.Println("Done.")
				os.Exit(1)
			}
		}
	}()

	if err := run(ctx, conf, task); err != nil {
		ErrorLog.Printf("%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, conf Config, task func()) error {
	conf.Init(os.Args)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			task()
		}
	}
}
