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
	IsInitiated() bool
}

type Daemon struct {
	conf  Config
	tasks []Task
}

func Create(conf Config) *Daemon {
	dmn := &Daemon{conf: conf}
	return dmn
}

func (d *Daemon) AddTask(task Task) *Daemon {
	d.tasks = append(d.tasks, task)
	return d
}

func (d *Daemon) AddTasks(tasks ...Task) *Daemon {
	d.tasks = append(d.tasks, tasks...)
	return d
}

func (d *Daemon) Start(ctx context.Context) (err error) {
	defer CatchWithReturn(&err)

	if !d.conf.IsInitiated() {
		err := d.conf.Init(os.Args)
		Throw(err)
	}

	ctx, cancel := context.WithCancel(ctx)

	go d.handleSignals(ctx, cancel)

	d.startTasks(ctx)
	return
}

func (d *Daemon) GetConfig() Config {
	return d.conf
}

func (d *Daemon) startTasks(ctx context.Context) {
	for _, t := range d.tasks {
		go t.Do(ctx)
	}
}

func (d *Daemon) handleSignals(ctx context.Context, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP:
			cancel()
			err := d.Start(ctx)
			Throw(err)
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
