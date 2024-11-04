package main

import (
	"context"
	"flag"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"os"
	"wombat/internal/app"
	"wombat/internal/dao"
	"wombat/pkg/daemon"
)

func main() {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())
	defer mainCancelCauseFunc(nil)

	conf := new(app.Config)
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
	kafkaHelper, err := app.NewKafkaHelper(kafkaConf)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	messageEventRepository, err := dao.NewMessageEventRepository(&conf.PostgreSQL.Url)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	telegram, err := app.NewTelegramSource(conf.Telegram.Token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	jiraHelper, err := app.NewJiraClient(conf.Jira.Url, conf.Jira.Username, conf.Jira.Password)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	dmn := daemon.Create(conf)

	runner, err := app.NewApplication(dmn, kafkaHelper, jiraHelper, messageEventRepository, telegram)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	runner.Run(mainCtx)
}

func parseArgs(args []string) ([]string, error) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}

	return []string{*configPath}, nil
}
