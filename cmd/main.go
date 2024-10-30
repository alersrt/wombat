package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"os"
	"wombat/internal/app"
	"wombat/internal/config"
	"wombat/internal/messaging"
	"wombat/internal/source"
	"wombat/pkg/daemon"
)

func main() {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())

	conf := new(config.Config)
	err := conf.Init(os.Args)
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

	updates := make(chan any)

	telegram := source.NewTelegramSource(updates, conf.Telegram.Token)

	dmn := daemon.Create(mainCtx, mainCancelCauseFunc, conf)

	app.NewApplication(
		mainCtx,
		mainCancelCauseFunc,
		updates,
		conf,
		kafkaHelper,
		telegram,
	).Run(dmn)
}
