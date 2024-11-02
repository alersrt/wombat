package daemon

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Daemon struct {
	conf  Config
	tasks []Task
}

func Create(conf Config) *Daemon {
	dmn := &Daemon{conf: conf}
	return dmn
}

func (receiver *Daemon) handleSignals(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	defer signal.Stop(signalChan)

	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP:
			err := receiver.conf.Init(os.Args)
			if err != nil {
				slog.Error(err.Error())
			}
		case os.Interrupt:
			os.Exit(130)
		case os.Kill:
			os.Exit(137)
		case syscall.SIGTERM:
			os.Exit(143)
		}
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		} else {
			slog.Info("Done.")
			os.Exit(0)
		}
	}
}

func (receiver *Daemon) AddTask(task Task) *Daemon {
	receiver.tasks = append(receiver.tasks, task)
	return receiver
}

func (receiver *Daemon) AddTasks(tasks ...Task) *Daemon {
	receiver.tasks = append(receiver.tasks, tasks...)
	return receiver
}

func (receiver *Daemon) Start(ctx context.Context) {
	if !receiver.conf.IsInitiated() {
		err := receiver.conf.Init(os.Args)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}

	go receiver.handleSignals(ctx)

	for _, task := range receiver.tasks {
		go task()
	}
}

func (receiver *Daemon) GetConfig() Config {
	return receiver.conf
}
