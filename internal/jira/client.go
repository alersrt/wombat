package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
)

type Client struct {
	client *jira.Client
}

func NewClient(url string, token string) (*Client, error) {
	tp := jira.PATAuthTransport{Token: token}
	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		return nil, fmt.Errorf("jira: client: new: %v", err)
	}
	return &Client{client}, nil
}

func (c *Client) Update(issue string, commentId string, text string) error {
	_, _, err := c.client.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	if err != nil {
		return fmt.Errorf("jira: client: update: %v", err)
	}
	return nil
}

func (c *Client) Add(issue string, text string) (string, error) {
	comment, _, err := c.client.Issue.AddComment(issue, &jira.Comment{Body: text})
	if err != nil {
		return "", fmt.Errorf("jira: client: add: %v", err)
	}
	return comment.ID, nil
}
