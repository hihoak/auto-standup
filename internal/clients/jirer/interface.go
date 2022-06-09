package jirer

import (
	"github.com/andygrunwald/go-jira"
)

//go:generate mockgen -destination "./mocks/mock_jirer.go" -package "mocks" -source "./interface.go" Jirer

// Jirer - interface for JiraClient
type Jirer interface {
	SearchIssue(jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error)
	GetIssue(key string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error)
}
