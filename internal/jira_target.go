package internal

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log/slog"
	"regexp"
	"wombat/internal/domain"
	"wombat/internal/storage"
	"wombat/pkg/cipher"
)

type TargetClient interface {
	Add(tag string, text string) (string, error)
	Update(tag string, commentId string, text string) error
}

type JiraClient struct {
	client *jira.Client
}

func NewJiraClient(url string, token string) (TargetClient, error) {
	tp := jira.PATAuthTransport{Token: token}
	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		return nil, fmt.Errorf("jira:client:new:err: url=%s", url)
	}
	return &JiraClient{client}, nil
}

func (c *JiraClient) Update(issue string, commentId string, text string) error {
	_, _, err := c.client.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	if err != nil {
		return fmt.Errorf("jira:client:update:err: %w", err)
	}
	return nil
}

func (c *JiraClient) Add(issue string, text string) (string, error) {
	comment, _, err := c.client.Issue.AddComment(issue, &jira.Comment{Body: text})
	if err != nil {
		return "", fmt.Errorf("jira:client:add:err: %w", err)
	}
	return comment.ID, nil
}

type JiraTarget struct {
	cipher     *cipher.AesGcmCipher
	targetType domain.TargetType
	url        string
	db         *storage.DbStorage
	tagsRegex  *regexp.Regexp
	router     *Router
}

func NewJiraTarget(
	url string,
	tag string,
	router *Router,
	dbStorage *storage.DbStorage,
	cipher *cipher.AesGcmCipher,
) *JiraTarget {
	return &JiraTarget{
		targetType: domain.TargetTypeJira,
		cipher:     cipher,
		url:        url,
		tagsRegex:  regexp.MustCompile(tag),
		db:         dbStorage,
		router:     router,
	}
}

func (t *JiraTarget) GetTargetType() domain.TargetType {
	return t.targetType
}

func (t *JiraTarget) Do(ctx context.Context) {
	slog.Info("jira:do:start")
	defer slog.Info("jira:do:finish")

	for {
		select {
		case req := <-t.router.ReqChan():
			err := t.handle(ctx, req)
			if err != nil {
				slog.Error(err.Error())
				t.router.SendRes(req.ToResponse(false, err.Error()))
			} else {
				t.router.SendRes(req.ToResponse(true, ""))
			}
		case <-ctx.Done():
			slog.Info("jira:do:ctx:done")
			return
		}
	}
}

func (t *JiraTarget) handle(ctx context.Context, req *domain.Request) error {
	ctxTx, cancelTx := context.WithCancel(ctx)
	defer cancelTx()

	if !t.tagsRegex.MatchString(req.Content) {
		return fmt.Errorf("jira:do:err: tag not found")
	}

	tx, err := t.db.BeginTx(ctxTx)
	if err != nil {
		return err
	}

	targetConnection, err := tx.GetTargetConnection(string(req.SourceType), string(req.TargetType), req.UserId)
	if err != nil {
		return err
	}
	token, err := t.cipher.Decrypt(targetConnection.Token)
	if err != nil {
		return err
	}
	client, err := NewJiraClient(t.url, token)
	if err != nil {
		return err
	}

	tags := t.tagsRegex.FindAllString(req.Content, -1)

	savedComments, err := tx.GetCommentMetadata(string(req.SourceType), req.ChatId, req.MessageId)
	if err != nil {
		return err
	}

	if len(savedComments) == 0 {
		for _, tag := range tags {
			commentId, err := client.Add(tag, req.Content)
			if err != nil {
				return err
			}
			_, err = tx.SaveCommentMetadata(&domain.Comment{Request: req, Tag: tag, CommentId: commentId})
			if err != nil {
				return err
			}
		}
	} else {
		taggedComments := map[string]*domain.Comment{}
		for _, comment := range savedComments {
			taggedComments[comment.Tag] = comment
		}
		for _, tag := range tags {
			comment := taggedComments[tag]
			err := client.Update(tag, comment.CommentId, req.Content)
			if err != nil {
				return err
			}
			_, err = tx.SaveCommentMetadata(&domain.Comment{Request: req, Tag: tag, CommentId: comment.CommentId})
			if err != nil {
				return err
			}
		}
	}
	if err = tx.CommitTx(); err != nil {
		return err
	}
	return nil
}
