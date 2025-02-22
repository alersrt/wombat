package app

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"os"
	"wombat/internal/domain"
)

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

		jiraHelper, err := NewJiraClient(receiver.conf.Jira.Url, conf.Jira.Username, conf.Jira.Token)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		tags := receiver.tagsRegex.FindAllString(msg.Text, -1)
		savedComments := receiver.commentRepository.GetMessagesByMetadata(msg.ChatId, msg.MessageId)
		if len(savedComments) == 0 {
			for _, tag := range tags {
				commentId, err := receiver.jiraHelper.AddComment(tag, msg.Text)
				if err != nil {
					return err
				}

				saved := receiver.commentRepository.Save(&domain.Comment{
					Message:   msg,
					Tag:       tag,
					CommentId: commentId,
				})

				slog.Info(fmt.Sprintf("Comsumed: %+v %+v", saved, saved.Message))
			}
		} else {
			taggedComments := map[string]*domain.Comment{}
			for _, comment := range savedComments {
				taggedComments[comment.Tag] = comment
			}
			for _, tag := range tags {
				comment := taggedComments[tag]
				err := receiver.jiraHelper.UpdateComment(tag, comment.CommentId, msg.Text)
				if err != nil {
					return err
				}

				saved := receiver.commentRepository.Save(&domain.Comment{
					Message:   msg,
					Tag:       tag,
					CommentId: comment.CommentId,
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
