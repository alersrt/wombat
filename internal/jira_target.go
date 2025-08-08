package internal

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"regexp"
	"wombat/pkg"
)

type TargetClient interface {
	Add(tag string, text string) string
	Update(tag string, commentId string, text string)
}

type JiraClient struct {
	*jira.Client
}

func NewJiraClient(url string, token string) (tc TargetClient, err error) {
	defer pkg.CatchWithReturn(&err)
	tp := jira.PATAuthTransport{Token: token}
	client, ex := jira.NewClient(tp.Client(), url)
	pkg.Throw(ex)
	if client != nil {
		tc = &JiraClient{client}
	}
	return
}

func (c *JiraClient) Update(issue string, commentId string, text string) {
	_, _, err := c.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	pkg.Throw(err)
	return
}

func (c *JiraClient) Add(issue string, text string) string {
	comment, _, ex := c.Issue.AddComment(issue, &jira.Comment{Body: text})
	pkg.Throw(ex)
	return comment.ID
}

type JiraTarget struct {
	cipher     *AesGcmCipher
	targetType TargetType
	url        string
	db         *DbStorage
	tagsRegex  *regexp.Regexp
	router     *Router
}

func NewJiraTarget(
	url string,
	tag string,
	router *Router,
	dbStorage *DbStorage,
	cipher *AesGcmCipher,
) *JiraTarget {
	return &JiraTarget{
		targetType: JiraType,
		cipher:     cipher,
		url:        url,
		tagsRegex:  regexp.MustCompile(tag),
		db:         dbStorage,
		router:     router,
	}
}

func (t *JiraTarget) GetTargetType() TargetType {
	return t.targetType
}

func (t *JiraTarget) Do(ctx context.Context) (err error) {
	defer pkg.CatchWithReturn(&err)
	for req := range t.router.GetReq() {
		ex := t.handle(ctx, req)
		res := &Response{
			SourceType: req.SourceType,
			TargetType: req.TargetType,
			UserId:     req.UserId,
			ChatId:     req.ChatId,
			MessageId:  req.MessageId,
		}
		if ex != nil {
			res.Ok = false
			t.router.SendRes(res)
			pkg.Throw(ex)
		} else {
			res.Ok = true
			t.router.SendRes(res)
		}
	}
	return
}

func (t *JiraTarget) handle(ctx context.Context, req *Request) (err error) {
	if !t.tagsRegex.MatchString(req.Content) {
		return fmt.Errorf("tag not found: %v", req.Content)
	}

	tx := t.db.BeginTx(ctx)
	defer pkg.CatchWithReturnAndCall(&err, tx.RollbackTx)

	targetConnection := tx.GetTargetConnection(req.SourceType.String(), req.TargetType.String(), req.UserId)
	client, ex := NewJiraClient(t.url, t.cipher.Decrypt(targetConnection.Token))
	pkg.Throw(ex)

	tags := t.tagsRegex.FindAllString(req.Content, -1)
	savedComments := tx.GetCommentMetadata(req.SourceType.String(), req.ChatId, req.MessageId)

	if len(savedComments) == 0 {
		for _, tag := range tags {
			commentId := client.Add(tag, req.Content)
			tx.SaveCommentMetadata(&Comment{Request: req, Tag: tag, CommentId: commentId})
		}
	} else {
		taggedComments := map[string]*Comment{}
		for _, comment := range savedComments {
			taggedComments[comment.Tag] = comment
		}
		for _, tag := range tags {
			comment := taggedComments[tag]
			client.Update(tag, comment.CommentId, req.Content)
			tx.SaveCommentMetadata(&Comment{Request: req, Tag: tag, CommentId: comment.CommentId})
		}
	}
	tx.CommitTx()
	return
}
