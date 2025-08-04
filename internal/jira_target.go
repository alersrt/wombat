package internal

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log/slog"
	"regexp"
)

type TargetClient interface {
	Add(tag string, text string) string
	Update(tag string, commentId string, text string)
}

type JiraClient struct {
	*jira.Client
}

func NewJiraClient(url string, token string) TargetClient {
	tp := jira.PATAuthTransport{
		Token: token,
	}
	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		panic(err)
	}
	return &JiraClient{client}
}

func (c *JiraClient) Update(issue string, commentId string, text string) {
	_, _, err := c.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	if err != nil {
		panic(err)
	}
	return
}

func (c *JiraClient) Add(issue string, text string) string {
	comment, _, err := c.Issue.AddComment(issue, &jira.Comment{
		Body: text,
	})
	if err != nil {
		panic(err)
	}
	return comment.ID
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
		t.handle(update)
	}
}

func (t *JiraTarget) handle(msg *Message) {
	if !t.tagsRegex.MatchString(msg.Content) {
		slog.Info(fmt.Sprintf("Tag not found: %v", msg.Content))
		return
	}

	tx := t.db.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.RollbackTx()
		}
	}()

	targetConnection := tx.GetTargetConnection(msg.SourceType.String(), msg.TargetType.String(), msg.UserId)
	client := NewJiraClient(t.url, targetConnection.Token)

	tags := t.tagsRegex.FindAllString(msg.Content, -1)
	savedComments := tx.GetCommentMetadata(msg.SourceType.String(), msg.ChatId, msg.MessageId)

	if len(savedComments) == 0 {
		for _, tag := range tags {
			commentId := client.Add(tag, msg.Content)
			tx.SaveCommentMetadata(&Comment{
				Message:   msg,
				Tag:       tag,
				CommentId: commentId,
			})
		}
	} else {
		taggedComments := map[string]*Comment{}
		for _, comment := range savedComments {
			taggedComments[comment.Tag] = comment
		}
		for _, tag := range tags {
			comment := taggedComments[tag]
			client.Update(tag, comment.CommentId, msg.Content)
			tx.SaveCommentMetadata(&Comment{
				Message:   msg,
				Tag:       tag,
				CommentId: comment.CommentId,
			})
		}
	}
	tx.CommitTx()
}
