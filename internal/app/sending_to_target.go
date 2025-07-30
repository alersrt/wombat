package app

import (
	"log/slog"
	"wombat/internal/domain"
)

func (receiver *Application) processTarget() {
	for update := range receiver.targetChan {
		tx, err := receiver.dbStorage.BeginTx()
		processError(err)

		client, err := NewJiraClient(receiver.conf.Jira.Url, "")
		processError(err)

		tags := receiver.tagsRegex.FindAllString(update.Content, -1)
		savedComments, err := tx.GetCommentsByMetadata(update.SourceType, update.ChatId, update.MessageId)
		processError(err)

		if len(savedComments) == 0 {
			for _, tag := range tags {
				commentId, err := client.Add(tag, update.Content)
				processError(err)
				tx.SaveComment(&domain.Comment{
					Message:   update,
					Tag:       tag,
					CommentId: commentId,
				})
			}
		} else {
			taggedComments := map[string]*domain.Comment{}
			for _, comment := range savedComments {
				taggedComments[comment.Tag] = comment
			}
			for _, tag := range tags {
				comment := taggedComments[tag]
				err := client.Update(tag, comment.CommentId, update.Content)
				processError(err)
				tx.SaveComment(&domain.Comment{
					Message:   update,
					Tag:       tag,
					CommentId: comment.CommentId,
				})
			}
		}
	}
}

func processError(err error) {
	if err != nil {
		slog.Warn(err.Error())
		return
	}
}
