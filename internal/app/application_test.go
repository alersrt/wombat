package app

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
	"wombat/internal/config"
	"wombat/internal/dao"
	"wombat/internal/domain"
	"wombat/internal/messaging"
	"wombat/internal/source"
	"wombat/pkg/daemon"
	"wombat/pkg/utils"
)

var (
	mockUpdatesChan = make(chan any)
)

var (
	messageEventRepository dao.QueryHelper[domain.MessageEvent, string]
)

func setup(
	ctx context.Context,
	kafkaContainer *testcontainers.DockerContainer,
	postgresContainer *testcontainers.DockerContainer,
) (*Application, error) {
	kafkaBootstrap, err := kafkaContainer.PortEndpoint(ctx, "9092", "")
	if err != nil {
		return nil, err
	}

	pgHost, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	pgPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	conf := &config.MockConfig{Config: &config.Config{
		Database: &config.Database{
			PostgreSQL: &config.PostgreSQL{
				Url: fmt.Sprintf("postgres://wombat_rw:wombat_rw@%s:%d/wombatdb?sslmode=disable", pgHost, pgPort.Int()),
			},
		},
		Kafka: &config.Kafka{
			GroupId:   "wombat",
			Bootstrap: kafkaBootstrap,
			Topic:     "wombat.test",
		},
		Bot: &config.Bot{Tag: "(TEST-\\d+)", Emoji: "👍"},
	}}
	err = conf.Init(os.Args)
	if err != nil {
		return nil, err
	}

	dmn := daemon.Create(conf.Config)

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.Bootstrap,
		"group.id":          conf.Kafka.GroupId,
		"auto.offset.reset": "earliest",
	}
	kafkaHelper, err := messaging.NewKafkaHelper(kafkaConf)
	if err != nil {
		return nil, err
	}

	messageEventRepository, err = dao.NewMessageEventRepository(&conf.PostgreSQL.Url)
	if err != nil {
		return nil, err
	}

	telegram := &source.MockSource{Source: mockUpdatesChan}

	return NewApplication(dmn, make(chan any), kafkaHelper, messageEventRepository, telegram)
}

func TestApplication(t *testing.T) {
	/*------ Arranges ------*/
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

	err = environment.
		WaitForService("flyway", wait.ForExit().WithPollInterval(1*time.Second)).
		WaitForService("kafka", wait.ForHealthCheck().WithPollInterval(1*time.Second)).
		Up(testCtx)
	require.NoError(t, environment.Up(testCtx), "compose.Up()")

	kafkaContainer, err := environment.ServiceContainer(testCtx, "kafka")
	require.NoError(t, err, "Kafka container")
	postgresContainer, err := environment.ServiceContainer(testCtx, "postgres")
	require.NoError(t, err, "PostgreSQL container")

	// Wait until `doneChan` is closed.
	testedUnit, err := setup(testCtx, kafkaContainer, postgresContainer)
	require.NoError(t, err, "setup()")

	succ := make(chan bool, 1)

	/*------ Actions ------*/
	go testedUnit.Run(testCtx)

	mockUpdatesChan <- tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text:      "TEST-100",
			From:      &tgbotapi.User{},
			Chat:      tgbotapi.Chat{ID: 1},
			MessageID: 1,
		},
	}

	/*------ Asserts ------*/
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
			}

			hash := uuid.NewSHA1(uuid.NameSpaceURL, []byte("TELEGRAM"+"1"+"1")).String()
			saved, geterr := messageEventRepository.GetById(hash)
			if geterr != nil || saved == nil {
				continue
			}

			if saved.Hash == hash {
				succ <- true
			}
		}
	}()

	select {
	case <-succ:
	case <-time.After(10 * time.Second):
		t.FailNow()
	}
}
