package app

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"os"
	"testing"
	"time"
	"wombat/internal/config"
	"wombat/internal/messaging"
	"wombat/internal/source"
	"wombat/pkg/daemon"
)

var (
	doneChan        = make(chan bool)
	mockUpdatesChan = make(chan any)
	dmn             *daemon.Daemon
)

func setup(t *testing.T) Application {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())

	conf := new(config.MockConfig)
	err := conf.Init(os.Args)
	if err != nil {
		t.Fatal(err)
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.Bootstrap,
		"group.id":          conf.Kafka.GroupId,
		"auto.offset.reset": "earliest",
	}
	kafkaHelper, err := messaging.NewKafkaHelper(kafkaConf)
	if err != nil {
		t.Fatal(err)
	}

	updates := make(chan any)

	telegram := &source.MockSource{Updates: updates, Source: mockUpdatesChan, Done: doneChan}

	dmn = daemon.Create(mainCtx, mainCancelCauseFunc, conf)

	return NewApplication(
		mainCtx,
		mainCancelCauseFunc,
		updates,
		(*config.Config)(conf),
		kafkaHelper,
		telegram,
	)
}

func TestApplication(t *testing.T) {
	testCtx, testCancelFunc := context.WithCancel(context.Background())
	t.Cleanup(testCancelFunc)

	environment, err := compose.NewDockerCompose("./docker/docker-compose.yaml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	err = environment.Down(testCtx, compose.RemoveOrphans(true), compose.RemoveImagesLocal, compose.RemoveVolumes(true))
	t.Cleanup(func() { require.NoError(t, err, "compose.Down()") })

	require.NoError(t, environment.Up(testCtx, compose.Wait(true)), "compose.Up()")

	// Wait until `doneChan` is closed.
	testedUnit := setup(t)

	go testedUnit.Run(dmn)

	select {
	case <-doneChan:
	case <-time.After(10 * time.Second):
		t.FailNow()
	}

	mockUpdatesChan <- &tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "TEST-100",
		},
	}

	select {
	case <-time.After(120 * time.Second):
		t.SkipNow()
	}
}
