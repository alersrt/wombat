package pkg

import (
	"context"
	"github.com/pkg/errors"
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

func (d *Daemon) Start(ctx context.Context) (int, error) {
	if !d.conf.IsInitiated() {
		err := d.conf.Init(os.Args)
		if err != nil {
			return 1, errors.New(err.Error())
		}
	}

	ctx, cancel := context.WithCancel(ctx)

	d.startTasks(ctx)

	return d.handleSignals(ctx, cancel)
}

func (d *Daemon) GetConfig() Config {
	return d.conf
}

func (d *Daemon) startTasks(ctx context.Context) {
	for _, t := range d.tasks {
		go t.Do(ctx)
	}
}

func (d *Daemon) handleSignals(ctx context.Context, cancel context.CancelFunc) (int, error) {
	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	for {
		select {
		case s := <-signalChan:
			switch s {
			case syscall.SIGHUP:
				cancel()
				code, err := d.Start(ctx)
				if err != nil {
					return code, err
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
