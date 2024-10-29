package app

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"log/slog"
	"os"
	"testing"
	"time"
	"wombat/internal/config"
	"wombat/internal/messaging"
)

var (
	done        = make(chan bool)
	mockUpdates = make(chan any)
)

type MockSource struct {
	Updates chan any
}

func (receiver *MockSource) Read(causeFunc context.CancelCauseFunc) {
	for update := range mockUpdates {
		receiver.Updates <- update
	}
}

func setupAndRun(t *testing.T) {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())

	conf := new(config.Config)
	err := conf.Init(os.Args)
	if err != nil {
		slog.Error(err.Error())
		t.Fatal(err)
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.Bootstrap,
		"group.id":          conf.Kafka.GroupId,
		"auto.offset.reset": "earliest",
	}
	kafkaHelper, err := messaging.NewKafkaHelper(kafkaConf)
	if err != nil {
		slog.Error(err.Error())
		t.Fatal(err)
	}

	updates := make(chan any)

	telegram := &MockSource{
		updates,
	}

	NewApplication(
		mainCtx,
		mainCancelCauseFunc,
		updates,
		conf,
		kafkaHelper,
		telegram,
	).Run()
	done <- true
}

func Test(t *testing.T) {
	testCtx, testCancelFunc := context.WithCancel(context.Background())

	environment, err := compose.NewDockerCompose("../../docker/docker-compose.yaml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(
			t,
			environment.Down(testCtx, compose.RemoveOrphans(true), compose.RemoveImagesLocal, compose.RemoveVolumes(true)),
			"compose.Down()",
		)
	})

	t.Cleanup(testCancelFunc)

	require.NoError(t, environment.Up(testCtx, compose.Wait(true)), "compose.Up()")

	// Wait until `done` is closed.
	setupAndRun(t)

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.FailNow()
	}

}
