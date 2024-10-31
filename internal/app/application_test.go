package app

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"os"
	"testing"
	"time"
	"wombat/internal/config"
	"wombat/internal/messaging"
	"wombat/internal/source"
	"wombat/pkg/daemon"
	"wombat/pkg/utils"
)

var (
	doneChan        = make(chan bool)
	mockUpdatesChan = make(chan any)
)

func setup(t *testing.T, kafkaContainer *testcontainers.DockerContainer) (Application, error) {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())

	kafkaBootstrap, err := kafkaContainer.PortEndpoint(mainCtx, "9092", "")
	if err != nil {
		t.Fatal(err)
	}

	conf := &config.MockConfig{Config: &config.Config{
		Kafka: &config.Kafka{
			GroupId:   "wombat",
			Bootstrap: kafkaBootstrap,
			Topic:     "wombat.test",
		},
		Bot: &config.Bot{Tag: "(TEST-\\d+)", Emoji: "👍"},
	}}
	err = conf.Init(os.Args)
	if err != nil {
		t.Fatal(err)
	}

	dmn := daemon.Create(mainCtx, mainCancelCauseFunc, conf.Config)

	updates := make(chan any)

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.Bootstrap,
		"group.id":          conf.Kafka.GroupId,
		"auto.offset.reset": "earliest",
	}
	kafkaHelper, err := messaging.NewKafkaHelper(kafkaConf)
	if err != nil {
		t.Fatal(err)
	}

	telegram := &source.MockSource{Updates: updates, Source: mockUpdatesChan, Done: doneChan}

	return NewApplication(
		dmn,
		updates,
		kafkaHelper,
		telegram,
	)
}

func TestApplication(t *testing.T) {
	testCtx, testCancelFunc := context.WithCancel(context.Background())
	t.Cleanup(testCancelFunc)

	composePath, err := utils.FindFilePath("docker", "docker-compose.yaml")
	require.NoError(t, err, "Compose location")
	environment, err := compose.NewDockerCompose(composePath)
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		err = environment.Down(testCtx, compose.RemoveOrphans(true), compose.RemoveVolumes(true), compose.RemoveImagesLocal)
		require.NoError(t, err, "compose.Down()")
	})

	require.NoError(t, environment.Up(testCtx, compose.Wait(true)), "compose.Up()")

	kafkaContainer, err := environment.ServiceContainer(testCtx, "kafka")
	require.NoError(t, err, "Kafka container")

	// Wait until `doneChan` is closed.
	testedUnit, err := setup(t, kafkaContainer)
	require.NoError(t, err, "setup()")

	go testedUnit.Run()

	select {
	case <-doneChan:
	case <-time.After(10 * time.Second):
		t.FailNow()
	}

	mockUpdatesChan <- tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text:      "TEST-100",
			From:      &tgbotapi.User{},
			Chat:      tgbotapi.Chat{ID: 1},
			MessageID: 1,
		},
	}

	select {
	case <-time.After(10 * time.Second):
	}
}
