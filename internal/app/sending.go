package app

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"wombat/internal/domain"
)

func (receiver *Application) getTargetClient(targetType domain.TargetType, sourceType domain.SourceType, authorId string) (TargetClient, error) {
	switch targetType {
	case domain.JIRA:
		return NewJiraClient(receiver.conf.Jira.Url, jiraToken)
	}
	return nil, nil
}

func (receiver *Application) send() {
	err := receiver.kafkaHelper.Subscribe([]string{receiver.conf.Kafka.Topic}, func(event *kafka.Message) error {
		msg := &domain.Message{}
		if err := json.Unmarshal(event.Value, msg); err != nil {
			return err
		}

		if !receiver.aclRepository.IsAuthorAllowed(msg.AuthorId) {
			slog.Info(fmt.Sprintf("Author is not allowed: %s", msg.AuthorId))
			return nil
		}

		client, err := receiver.getTargetClient(domain.JIRA, msg.SourceType, msg.AuthorId)
		if err != nil {
			return err
		}

		tags := receiver.tagsRegex.FindAllString(msg.Text, -1)
		savedComments := receiver.commentRepository.GetMessagesByMetadata(msg.ChatId, msg.MessageId)
		if len(savedComments) == 0 {
			for _, tag := range tags {
				commentId, err := client.Add(tag, msg.Text)
				if err != nil {
					return err
				}

				saved := receiver.commentRepository.Save(&domain.Comment{
					TargetType: domain.JIRA,
					Message:    msg,
					Tag:        tag,
					CommentId:  commentId,
				})

				slog.Info(fmt.Sprintf("Consumed: %+v %+v", saved, saved.Message))
			}
		} else {
			taggedComments := map[string]*domain.Comment{}
			for _, comment := range savedComments {
				taggedComments[comment.Tag] = comment
			}
			for _, tag := range tags {
				comment := taggedComments[tag]
				err := client.Update(tag, comment.CommentId, msg.Text)
				if err != nil {
					return err
				}

				saved := receiver.commentRepository.Save(&domain.Comment{
					TargetType: domain.JIRA,
					Message:    msg,
					Tag:        tag,
					CommentId:  comment.CommentId,
				})

				slog.Info(fmt.Sprintf("Comsumed: %+v", saved))
			}
		}

		return nil
	})

	if err != nil {
		slog.Warn(err.Error())
	}
}
