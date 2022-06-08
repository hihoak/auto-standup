package jirer

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/pkg/utils"
)

type Client struct {
	Client                 *jira.Client
}

func New(username, password string) (*Client, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := jira.NewClient(tp.Client(), "https://jit.ozon.ru")
	if err != nil {
		utils.Log.Debug().Err(err).Msg("Can't initialize Jira client!")
		return nil, fmt.Errorf("can't initialize Jira client. Error: %w", err)
	}

	return &Client{
		Client: client,
	}, nil
}

func (j *Client) GetIssue(key string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error) {
	utils.Log.Debug().Msgf("Trying GetIssue with key '%s' and options: %v", key, options)
	issue, resp, err := j.Client.Issue.Get(key, options)
	utils.Log.Debug().Msgf("Result of GetIssue: %v", issue)
	return issue, resp, err
}

func (j *Client) SearchIssue(jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error) {
	utils.Log.Debug().Msgf("Trying SearchIssue with JQL '%s' and options: %v", jql, options)
	issues, resp, err := j.Client.Issue.Search(jql, options)
	utils.Log.Debug().Msgf("Result of SearchIssue: %v", issues)
	return issues, resp, err
}
