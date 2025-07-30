package app

import (
	"fmt"
	"log/slog"
	"wombat/internal/domain"
)

func (receiver *Application) processTarget() {
	for update := range receiver.targetChan {
		isAllowed, err := receiver.dbStorage.IsAuthorAllowed(update.SourceType, update.AuthorId)
		processError(err)
		if !isAllowed {
			slog.Info(fmt.Sprintf("Author is not allowed: %s:%s", update.SourceType.String(), update.AuthorId))
			return
		}

		client, err := NewJiraClient(receiver.conf.Jira.Url, "")
		processError(err)
		tags := receiver.tagsRegex.FindAllString(update.Text, -1)
		savedComments, err := receiver.dbStorage.GetCommentsByMetadata(update.SourceType, update.ChatId, update.MessageId)
		processError(err)
		if len(savedComments) == 0 {
			for _, tag := range tags {
				commentId, err := client.Add(tag, update.Text)
				processError(err)
				receiver.dbStorage.SaveComment(&domain.Comment{
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
				err := client.Update(tag, comment.CommentId, update.Text)
				processError(err)
				receiver.dbStorage.SaveComment(&domain.Comment{
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
