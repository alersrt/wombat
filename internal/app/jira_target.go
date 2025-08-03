package app

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log/slog"
	"regexp"
	"wombat/internal/domain"
	"wombat/internal/storage"
)

type TargetClient interface {
	Add(tag string, text string) (string, error)
	Update(tag string, commentId string, text string) error
}

type JiraClient struct {
	*jira.Client
}

func NewJiraClient(url string, token string) (TargetClient, error) {
	tp := jira.PATAuthTransport{
		Token: token,
	}
	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		return nil, err
	}
	return &JiraClient{
		client,
	}, nil
}

func (receiver *JiraClient) Update(issue string, commentId string, text string) error {
	_, _, err := receiver.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	return err
}

func (receiver *JiraClient) Add(issue string, text string) (string, error) {
	comment, _, err := receiver.Issue.AddComment(issue, &jira.Comment{
		Body: text,
	})
	if err != nil {
		return "", err
	}
	return comment.ID, nil
}

type JiraTarget struct {
	targetType domain.TargetType
	url        string
	db         *storage.DbStorage
	tagsRegex  *regexp.Regexp
	srcChan    chan *domain.Message
}

func NewJiraTarget(
	url string,
	tag string,
	dbStorage *storage.DbStorage,
	srcChan chan *domain.Message,
) (*JiraTarget, error) {
	return &JiraTarget{
		url:       url,
		tagsRegex: regexp.MustCompile(tag),
		db:        dbStorage,
		srcChan:   srcChan,
	}, nil
}

func (receiver *JiraTarget) GetTargetType() domain.TargetType {
	return receiver.targetType
}

func (receiver *JiraTarget) Process() {
	for update := range receiver.srcChan {
		if !receiver.tagsRegex.MatchString(update.Content) {
			slog.Info(fmt.Sprintf("Tag not found: %v", update.Content))
			return
		}

		tx, err := receiver.db.BeginTx()
		processError(err, tx)

		targetConnection, err := tx.GetTargetConnection(update.SourceType.String(), update.TargetType.String(), update.UserId)
		processError(err, tx)

		client, err := NewJiraClient(receiver.url, targetConnection.Token)
		processError(err, tx)

		tags := receiver.tagsRegex.FindAllString(update.Content, -1)
		savedComments, err := tx.GetCommentMetadata(update.SourceType.String(), update.ChatId, update.MessageId)
		processError(err, tx)

		if len(savedComments) == 0 {
			for _, tag := range tags {
				commentId, err := client.Add(tag, update.Content)
				processError(err, tx)
				tx.SaveCommentMetadata(&domain.Comment{
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
				processError(err, tx)
				tx.SaveCommentMetadata(&domain.Comment{
					Message:   update,
					Tag:       tag,
					CommentId: comment.CommentId,
				})
			}
		}
	}
}

func processError(err error, tx *storage.Tx) {
	if err != nil {
		slog.Warn(err.Error())
		if tx != nil {
			tx.RollbackTx()
		}
		return
	}
}
