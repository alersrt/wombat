package internal

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log/slog"
	"regexp"
	"wombat/pkg"
)

type TargetClient interface {
	Add(tag string, text string) (string, error)
	Update(tag string, commentId string, text string) error
}

type JiraClient struct {
	*jira.Client
}

func NewJiraClient(url string, token string) (tc TargetClient, err error) {
	defer pkg.Catch(&err)
	tp := jira.PATAuthTransport{Token: token}
	client, ex := jira.NewClient(tp.Client(), url)
	pkg.Try(ex)
	if client != nil {
		tc = &JiraClient{client}
	}
	return
}

func (c *JiraClient) Update(issue string, commentId string, text string) (err error) {
	_, _, err = c.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	return
}

func (c *JiraClient) Add(issue string, text string) (cId string, err error) {
	defer pkg.Catch(&err)
	comment, _, ex := c.Issue.AddComment(issue, &jira.Comment{Body: text})
	pkg.Try(ex)
	if comment != nil {
		cId = comment.ID
	}
	return
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
) *JiraTarget {
	return &JiraTarget{
		url:        url,
		tagsRegex:  regexp.MustCompile(tag),
		db:         dbStorage,
		srcChan:    srcChan,
		targetType: JiraType,
	}
}

func (t *JiraTarget) GetTargetType() TargetType {
	return t.targetType
}

func (t *JiraTarget) Do(ctx context.Context) (err error) {
	defer pkg.Catch(&err)
	for update := range t.srcChan {
		ex := t.handle(update)
		pkg.Try(ex)
	}
	return
}

func (t *JiraTarget) handle(msg *Message) (err error) {
	if !t.tagsRegex.MatchString(msg.Content) {
		slog.Info(fmt.Sprintf("Tag not found: %v", msg.Content))
		return
	}

	tx := t.db.BeginTx()
	defer pkg.CatchWithPost(&err, tx.RollbackTx)

	targetConnection := tx.GetTargetConnection(msg.SourceType.String(), msg.TargetType.String(), msg.UserId)
	client, ex := NewJiraClient(t.url, targetConnection.Token)
	pkg.Try(ex)

	tags := t.tagsRegex.FindAllString(msg.Content, -1)
	savedComments := tx.GetCommentMetadata(msg.SourceType.String(), msg.ChatId, msg.MessageId)

	if len(savedComments) == 0 {
		for _, tag := range tags {
			commentId, ex := client.Add(tag, msg.Content)
			pkg.Try(ex)
			tx.SaveCommentMetadata(&Comment{Message: msg, Tag: tag, CommentId: commentId})
		}
	} else {
		taggedComments := map[string]*Comment{}
		for _, comment := range savedComments {
			taggedComments[comment.Tag] = comment
		}
		for _, tag := range tags {
			comment := taggedComments[tag]
			ex := client.Update(tag, comment.CommentId, msg.Content)
			pkg.Try(ex)
			tx.SaveCommentMetadata(&Comment{Message: msg, Tag: tag, CommentId: comment.CommentId})
		}
	}
	tx.CommitTx()
	return
}
