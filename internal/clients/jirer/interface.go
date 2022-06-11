package jirer

import (
	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/pkg/utils"
)

//go:generate mockgen -destination "./mocks/mock_jirer.go" -package "mocks" -source "./interface.go" Jirer

// Jirer - interface for JiraClient
type Jirer interface {
	SearchIssue(jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error)
	GetIssue(key string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error)
	GetIssueLogTimeForTheLastWorkDay(cfg *utils.Config, issue jira.Issue) int
}
