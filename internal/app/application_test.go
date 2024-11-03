package app

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
	"wombat/internal/dao"
	"wombat/internal/domain"
	"wombat/pkg/daemon"
)

type MockConfig struct {
	*Config
}

func (receiver *MockConfig) Init(args []string) error {
	receiver.isInitiated = true
	return nil
}

func (receiver *MockConfig) IsInitiated() bool {
	return true
}

type MockSource struct {
	SourceChan chan *domain.MessageEvent
}

func (receiver *MockSource) ForwardTo(target chan *domain.MessageEvent) {
	for update := range receiver.SourceChan {
		target <- update
	}
}

var (
	mockUpdatesChan = make(chan *domain.MessageEvent)
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

	conf := &MockConfig{Config: &Config{
		Database: &Database{
			PostgreSQL: &PostgreSQL{
				Url: fmt.Sprintf("postgres://wombat_rw:wombat_rw@%s:%d/wombatdb?sslmode=disable", pgHost, pgPort.Int()),
			},
		},
		Kafka: &Kafka{
			GroupId:   "wombat",
			Bootstrap: kafkaBootstrap,
			Topic:     "wombat.test",
		},
		Bot: &Bot{Tag: "(TEST-\\d+)", Emoji: "👍"},
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
	kafkaHelper, err := NewKafkaHelper(kafkaConf)
	if err != nil {
		return nil, err
	}

	messageEventRepository, err = dao.NewMessageEventRepository(&conf.PostgreSQL.Url)
	if err != nil {
		return nil, err
	}

	telegram := &MockSource{SourceChan: mockUpdatesChan}

	return NewApplication(dmn, kafkaHelper, messageEventRepository, telegram)
}

func TestApplication(t *testing.T) {
	/*------ Arranges ------*/
	testCtx, testCancelFunc := context.WithCancel(context.Background())
	t.Cleanup(testCancelFunc)

	composePath, err := FindFilePath("docker", "docker-compose.yaml")
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

	hash := uuid.NewSHA1(uuid.NameSpaceURL, []byte("TELEGRAM"+"1"+"1")).String()
	mockUpdatesChan <- &domain.MessageEvent{
		Hash:       hash,
		SourceType: domain.TELEGRAM,
		EventType:  domain.CREATE,
		ChatId:     "1",
		MessageId:  "1",
		Text:       "TEST-100",
	}

	/*------ Asserts ------*/
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
			}

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
