package jirer

import (
	"fmt"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/pkg/utils"
)

// Client - Jira Client
type Client struct {
	Client *jira.Client
}

// New - creator for Jira Client
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

// GetIssue - returning Jira Issue
func (j *Client) GetIssue(key string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error) {
	utils.Log.Debug().Msgf("Trying GetIssue with key '%s' and options: %v", key, options)
	issue, resp, err := j.Client.Issue.Get(key, options)
	utils.Log.Debug().Msgf("Result of GetIssue: %v", issue)
	return issue, resp, err
}

// SearchIssue - Searching for issues
func (j *Client) SearchIssue(jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error) {
	utils.Log.Debug().Msgf("Trying SearchIssue with JQL '%s' and options: %v", jql, options)
	issues, resp, err := j.Client.Issue.Search(jql, options)
	utils.Log.Debug().Msgf("Result of SearchIssue: %v", issues)
	return issues, resp, err
}

// GetIssueLogTimeForTheLastWorkDay - ...
func (j *Client) GetIssueLogTimeForTheLastWorkDay(cfg *utils.Config, issue jira.Issue) int {
	startTime := cfg.CmdStartTime.AddDate(0, 0, -cfg.NumberOfDaysForGetTickets)
	endTime := cfg.CmdStartTime
	totalLogTimeForLastWorkDay := 0
	if issue.Fields.Worklog == nil || issue.Fields.Worklog.Worklogs == nil {
		utils.Log.Debug().Msgf("Empty worklog of ticket %v", issue)
		return 0
	}
	for _, worklog := range issue.Fields.Worklog.Worklogs {
		logTime := time.Time(*worklog.Created).UTC()
		if startTime.Before(logTime) && endTime.After(logTime) {
			totalLogTimeForLastWorkDay += worklog.TimeSpentSeconds
		}
	}
	return totalLogTimeForLastWorkDay
}
