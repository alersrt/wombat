package main

import (
	"context"
	"flag"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"os"
	"wombat/internal/app"
	"wombat/internal/config"
	"wombat/internal/dao"
	"wombat/internal/messaging"
	"wombat/internal/source"
	"wombat/pkg/daemon"
)

func main() {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())

	conf := new(config.Config)
	args, err := parseArgs(os.Args)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	err = conf.Init(args)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.Bootstrap,
		"group.id":          conf.Kafka.GroupId,
		"auto.offset.reset": "earliest",
	}
	kafkaHelper, err := messaging.NewKafkaHelper(kafkaConf)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	pgUrl := conf.PostgreSQL.FormatURL()
	queryHelper := dao.NewPostgreSQLManager(mainCtx, &pgUrl)

	updates := make(chan any)

	telegram, err := source.NewTelegramSource(conf.Telegram.Token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	dmn := daemon.Create(mainCtx, mainCancelCauseFunc, conf)

	runner, err := app.NewApplication(dmn, updates, kafkaHelper, queryHelper, telegram)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	runner.Run()
}

func parseArgs(args []string) ([]string, error) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}

	return []string{*configPath}, nil
}
