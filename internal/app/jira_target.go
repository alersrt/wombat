package app

import "github.com/andygrunwald/go-jira"

type JiraClient struct {
	client *jira.Client
}

func NewJiraClient(url string, token string) (*JiraClient, error) {
	tp := jira.BearerAuthTransport{
		Token: token,
	}
	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		return nil, err
	}
	return &JiraClient{
		client: client,
	}, nil
}

func (receiver *JiraClient) Update(issue string, commentId string, text string) error {
	_, _, err := receiver.client.Issue.UpdateComment(issue, &jira.Comment{ID: commentId, Body: text})
	return err
}

func (receiver *JiraClient) Add(issue string, text string) (string, error) {
	comment, _, err := receiver.client.Issue.AddComment(issue, &jira.Comment{
		Body: text,
	})
	if err != nil {
		return "", err
	}
	return comment.ID, nil
}
