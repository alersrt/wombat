package jira

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"wombat/internal/domain"
	router2 "wombat/internal/router"
	"wombat/internal/storage"
	"wombat/pkg/cipher"
)

type TargetClient interface {
	Add(tag string, text string) (string, error)
	Update(tag string, commentId string, text string) error
}

type Target struct {
	cipher     *cipher.AesGcmCipher
	targetType domain.TargetType
	url        string
	db         *storage.DbStorage
	tagsRegex  *regexp.Regexp
	router     *router2.Router
}

func NewJiraTarget(
	url string,
	tag string,
	router *router2.Router,
	dbStorage *storage.DbStorage,
	cipher *cipher.AesGcmCipher,
) *Target {
	return &Target{
		targetType: domain.TargetTypeJira,
		cipher:     cipher,
		url:        url,
		tagsRegex:  regexp.MustCompile(tag),
		db:         dbStorage,
		router:     router,
	}
}

func (t *Target) Do(ctx context.Context) {
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

func (t *Target) handle(ctx context.Context, req *domain.Request) error {
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
	client, err := NewClient(t.url, token)
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
