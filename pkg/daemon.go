package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Task func(ctx context.Context)

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

func (d *Daemon) handleSignals(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	defer signal.Stop(signalChan)

	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP:
			err := d.conf.Init(os.Args)
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

func (d *Daemon) AddTask(task Task) *Daemon {
	d.tasks = append(d.tasks, task)
	return d
}

func (d *Daemon) AddTasks(tasks ...Task) *Daemon {
	d.tasks = append(d.tasks, tasks...)
	return d
}

func (d *Daemon) Start(ctx context.Context) {
	if !d.conf.IsInitiated() {
		err := d.conf.Init(os.Args)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}

	go d.handleSignals(ctx)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	for _, task := range d.tasks {
		go task(ctx)
	}
}

func (d *Daemon) GetConfig() Config {
	return d.conf
}
