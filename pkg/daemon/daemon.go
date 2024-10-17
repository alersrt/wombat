package daemon

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()

	if err := run(ctx, conf, task); err != nil {
		log.Printf("%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, conf Config, task func()) error {
	conf.Init(os.Args)
	log.SetOutput(os.Stdout)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			task()
		}
	}
}
