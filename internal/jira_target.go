package internal

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log/slog"
	"regexp"
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

func (c *JiraClient) Update(issue string, commentId string, text string) error {
	_, _, err := c.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	return err
}

func (c *JiraClient) Add(issue string, text string) (string, error) {
	comment, _, err := c.Issue.AddComment(issue, &jira.Comment{
		Body: text,
	})
	if err != nil {
		return "", err
	}
	return comment.ID, nil
}

type JiraTarget struct {
	targetType TargetType
	url        string
	db         *DbStorage
	tagsRegex  *regexp.Regexp
	srcChan    chan *Message
}

func NewJiraTarget(
	url string,
	tag string,
	dbStorage *DbStorage,
	srcChan chan *Message,
) (*JiraTarget, error) {
	return &JiraTarget{
		url:        url,
		tagsRegex:  regexp.MustCompile(tag),
		db:         dbStorage,
		srcChan:    srcChan,
		targetType: JiraType,
	}, nil
}

func (t *JiraTarget) GetTargetType() TargetType {
	return t.targetType
}

func (t *JiraTarget) Do(ctx context.Context) {
	for update := range t.srcChan {
		if !t.tagsRegex.MatchString(update.Content) {
			slog.Info(fmt.Sprintf("Tag not found: %v", update.Content))
			return
		}

		tx, err := t.db.BeginTx()
		processError(err, tx)

		targetConnection, err := tx.GetTargetConnection(update.SourceType.String(), update.TargetType.String(), update.UserId)
		processError(err, tx)

		client, err := NewJiraClient(t.url, targetConnection.Token)
		processError(err, tx)

		tags := t.tagsRegex.FindAllString(update.Content, -1)
		savedComments, err := tx.GetCommentMetadata(update.SourceType.String(), update.ChatId, update.MessageId)
		processError(err, tx)

		if len(savedComments) == 0 {
			for _, tag := range tags {
				commentId, err := client.Add(tag, update.Content)
				processError(err, tx)
				_, err = tx.SaveCommentMetadata(&Comment{
					Message:   update,
					Tag:       tag,
					CommentId: commentId,
				})
				processError(err, tx)
			}
		} else {
			taggedComments := map[string]*Comment{}
			for _, comment := range savedComments {
				taggedComments[comment.Tag] = comment
			}
			for _, tag := range tags {
				comment := taggedComments[tag]
				err := client.Update(tag, comment.CommentId, update.Content)
				processError(err, tx)
				_, err = tx.SaveCommentMetadata(&Comment{
					Message:   update,
					Tag:       tag,
					CommentId: comment.CommentId,
				})
				processError(err, tx)
			}
		}
		tx.CommitTx()
	}
}

func processError(err error, tx *Tx) {
	if err != nil {
		slog.Warn(err.Error())
		if tx != nil {
			tx.RollbackTx()
		}
		return
	}
}
