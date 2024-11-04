package app

import "github.com/andygrunwald/go-jira"

type JiraHelper interface {
	AddComment(issue string, text string) (string, error)
	UpdateComment(issue string, commentId string, text string) error
}

type JiraClient struct {
	client *jira.Client
}

func NewJiraClient(url string) (*JiraClient, error) {
	client, err := jira.NewClient(nil, url)
	if err != nil {
		return nil, err
	}
	return &JiraClient{
		client: client,
	}, nil
}

func (receiver *JiraClient) UpdateComment(issue string, commentId string, text string) error {
	_, _, err := receiver.client.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	return err
}

func (receiver *JiraClient) AddComment(issue string, text string) (string, error) {
	comment, _, err := receiver.client.Issue.AddComment(issue, &jira.Comment{
		Body: text,
	})
	if err != nil {
		return "", err
	}
	return comment.ID, nil
}
