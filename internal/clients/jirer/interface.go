package jirer

import (
	"github.com/andygrunwald/go-jira"
)

//go:generate mockgen -destination "./mock_jirer.go" -package "jirer" -source "./interface.go" Jirer

type Jirer interface {
	SearchIssue(jql string, options *jira.SearchOptions) ([]jira.Issue, *jira.Response, error)
	GetIssue(key string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error)
}
